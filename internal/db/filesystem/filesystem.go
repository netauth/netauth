// Package filesystem implements a key/value store on top of a generic
// filesystem.  This is the direct successor to the protodb and is
// compatible with its storage format.  It is not compatible with some
// of the features of the protodb, most notably noticing changes to
// the filesystem outside of NetAuth.  It was incredibly hard to make
// this work reliably in protodb, and if you look too closely you'll
// realize that it doesn't satisfy a lot of integrity constraints and
// probably could be used to corrupt data if you were really clever.
// Additionally, the filesystem key/value store does not use the .dat
// extension on data files as it is wholely unnecessary.  This needs
// to be done during migration.  The recommended way to migrate from
// one to another is to use a shell fragment that can talk to both.
package filesystem

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	atomic "github.com/google/renameio"
	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"

	"github.com/netauth/netauth/internal/db"
	"github.com/netauth/netauth/internal/startup"
)

// Filesystem anchors all the methods in the filesystem key/value
// store.
type Filesystem struct {
	basePath string

	l  hclog.Logger
	eF func(db.Event)
}

// event is an enum for what type of event to fire and subsequently
// map to a DB event.
type eventType int

const (
	eventUpdate eventType = iota
	eventDelete
)

var (
	// ErrPathEscape is returned if a key tries to climb up and
	// out of a directory.
	ErrPathEscape = errors.New("attempted path escape")
)

func init() {
	startup.RegisterCallback(cb)
}

func cb() {
	db.RegisterKV("filesystem", newKV)
}

func newKV(l hclog.Logger) (db.KVStore, error) {
	x := &Filesystem{
		l: l.Named("filesystem"),

		basePath: filepath.Join(viper.GetString("core.home"), "kv"),
	}

	return x, nil
}

// SetEventFunc sets up a function to call to fire events to
// subscribers.
func (fs *Filesystem) SetEventFunc(ef func(db.Event)) {
	fs.eF = ef
}

// Put stores a series of bytes on the filesystem, checking to make
// sure that the path is inside of the basePath
func (fs *Filesystem) Put(_ context.Context, k string, v []byte) error {
	p, err := fs.cleanPath(k)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(p), 0750); err != nil {
		return err
	}

	if err := atomic.WriteFile(p, v, 0640); err != nil {
		return err
	}

	fs.fireEventForKey(k, eventUpdate)
	return nil
}

// Get returns a series of bytes from the filesystem, checking to make
// sure that the bytes come from inside the base path.
func (fs *Filesystem) Get(_ context.Context, k string) ([]byte, error) {
	p, err := fs.cleanPath(k)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadFile(p)
	if os.IsNotExist(err) {
		return nil, db.ErrNoValue
	}
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// Del removes a file from disk that is inside the base path.
func (fs *Filesystem) Del(_ context.Context, k string) error {
	p, err := fs.cleanPath(k)
	if err != nil {
		return err
	}

	if _, err := os.Stat(p); os.IsNotExist(err) {
		return db.ErrNoValue
	}

	if err := os.Remove(p); err != nil {
		return err
	}

	fs.fireEventForKey(k, eventDelete)
	return nil
}

// Keys is a way to enumerate the keys in the key/value store and to
// optionally filter them based on a globbing expression.  This cheats
// and uses superior knowledge that NetAuth uses only a single key
// namespace with a single layer of keys below it.  Its technically
// possible to do something dumb with an entity or group name that
// includes a path seperator, but this should be filtered out at a
// higher level.
func (fs *Filesystem) Keys(_ context.Context, f string) ([]string, error) {
	// Discard error because the hard coded pattern cannot return
	// an os.PathError
	keys, _ := filepath.Glob(filepath.Join(fs.basePath, "*", "*"))

	out := make([]string, len(keys))
	i := 0
	for _, k := range keys {
		k, _ = filepath.Rel(fs.basePath, k)
		k = "/" + k
		if m, _ := filepath.Match(f, k); m {
			out[i] = k
			fs.l.Trace("Matched filter", "key", k)
			i++
		}
	}
	return out[:i], nil
}

// Close is required by the interface, but all operations on the
// filesystem are atomic, so no close is required.
func (fs *Filesystem) Close() error { return nil }

// Capabilities returns the capabilities that this implementation is
// able to satisfy.  Capabilities checks for a .writeable flag to tell
// it that the local copy is intentionally mutable.  Calls to Put may
// succeed even if this flag is missing, but higher level constructs
// can use this to check of this instance is in read-only mode.
func (fs *Filesystem) Capabilities() []db.KVCapability {
	out := []db.KVCapability{}

	// A file is used for this flag because it more natively
	// integrates with the keyspace concept of the filesystem.
	// Its also a clear example of putting the indicator in the
	// path that you are required to traverse.  An admin must be
	// aware of this file for their server to work.  This clues
	// them in that this file should not be replicated to
	// elsewhere, and what it does.
	_, err := os.Stat(filepath.Join(fs.basePath, ".mutable"))
	if !os.IsNotExist(err) {
		out = append(out, db.KVMutable)
	}

	return out
}

// cleanPath ensures that the path is inside of the base path.  This
// is only promised to work on *nix systems, as Windows is an unholy
// hellscape of legacy support that I'm pretty sure would let you
// escape a moderate quality safe given enough time.
func (fs *Filesystem) cleanPath(p string) (string, error) {
	p = filepath.Clean(p)

	if strings.HasPrefix(p, "..") {
		return "", ErrPathEscape
	}

	return filepath.Join(fs.basePath, p), nil
}

// fireEventForKey maps from a key to an entity or group and fires an
// appropriate event for the given key.
func (fs *Filesystem) fireEventForKey(k string, t eventType) {
	switch {
	case t == eventUpdate && strings.HasPrefix(k, "/entities/"):
		fs.eF(db.Event{
			PK:   filepath.Base(k),
			Type: db.EventEntityUpdate,
		})
	case t == eventDelete && strings.HasPrefix(k, "/entities/"):
		fs.eF(db.Event{
			PK:   filepath.Base(k),
			Type: db.EventEntityDestroy,
		})
	case t == eventUpdate && strings.HasPrefix(k, "/groups/"):
		fs.eF(db.Event{
			PK:   filepath.Base(k),
			Type: db.EventGroupUpdate,
		})
	case t == eventDelete && strings.HasPrefix(k, "/groups/"):
		fs.eF(db.Event{
			PK:   filepath.Base(k),
			Type: db.EventGroupDestroy,
		})
	default:
		fs.l.Warn("Event translation called with unknown key prefix", "type", t, "key", k)
	}
}
