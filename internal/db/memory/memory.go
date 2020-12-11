// Package memory implements a fully in-memory key/value store for the
// database to use.
package memory

import (
	"path"
	"strings"
	"sync"

	"github.com/hashicorp/go-hclog"

	"github.com/netauth/netauth/internal/db"
	"github.com/netauth/netauth/internal/startup"
)

// KV is a fully in-memory KV store.  It is intended to be used
// with integration tests, not in production.  It is not included in
// release builds by default.
type KV struct {
	sync.RWMutex
	m map[string][]byte
	l hclog.Logger

	eF func(db.Event)
}

func init() {
	startup.RegisterCallback(cb)
}

func cb() {
	db.RegisterKV("memory", NewKV)
}

// NewKV is exported to allow other test stores to be built on top of
// this one, such as the errorableKV used in the rpc2 integration
// test.
func NewKV(l hclog.Logger) (db.KVStore, error) {
	x := &KV{
		RWMutex: sync.RWMutex{},
		m:       make(map[string][]byte),
		l:       l.Named("memory"),
	}
	x.l.Debug("Initialization Complete")
	return x, nil
}

// SetEventFunc supplies the event firing entrypoing for this
// implementation.
func (kv *KV) SetEventFunc(ef func(db.Event)) {
	kv.eF = ef
}

// Put stores a value
func (kv *KV) Put(k string, v []byte) error {
	kv.Lock()
	kv.l.Trace("PUT", "key", k, "value", v)
	kv.m[k] = v
	kv.Unlock()

	switch {
	case strings.HasPrefix(k, "/entities"):
		kv.eF(db.Event{
			Type: db.EventEntityUpdate,
			PK:   path.Base(k),
		})
	case strings.HasPrefix(k, "/groups"):
		kv.eF(db.Event{
			Type: db.EventGroupUpdate,
			PK:   path.Base(k),
		})
	}
	return nil
}

// Get retrives a value
func (kv *KV) Get(k string) ([]byte, error) {
	kv.RLock()
	defer kv.RUnlock()
	v, ok := kv.m[k]
	if !ok {
		kv.l.Trace("NoValue", "key", k)
		return nil, db.ErrNoValue
	}
	kv.l.Trace("GET", "key", k, "value", v)
	return v, nil
}

// Del removes a value for a given key
func (kv *KV) Del(k string) error {
	kv.Lock()
	delete(kv.m, k)
	kv.Unlock()
	kv.l.Trace("DEL", "key", k)

	switch {
	case strings.HasPrefix(k, "/entities"):
		kv.eF(db.Event{
			Type: db.EventEntityDestroy,
			PK:   path.Base(k),
		})
	case strings.HasPrefix(k, "/groups"):
		kv.eF(db.Event{
			Type: db.EventGroupDestroy,
			PK:   path.Base(k),
		})
	}
	return nil
}

// Keys returns a set of keys optionally filtered by the filter
// expression.  filter should be a regex/shell type glob.
func (kv *KV) Keys(filter string) ([]string, error) {
	kv.RLock()
	defer kv.RUnlock()

	out := []string{}
	for k := range kv.m {
		// Only possible error here is ErrBadPattern, but all
		// the input patterns are hard-coded in other parts of
		// the code base.
		if match, _ := path.Match(filter, k); match {
			out = append(out, k)
		}
	}
	kv.l.Trace("KEYS", "filter", filter, "out", out)
	return out, nil
}

// Close is used in other implementations to finalize access.
func (kv *KV) Close() error { return nil }

// Capabilities is used to interrogate a KV store for capabilities.
func (kv *KV) Capabilities() []db.KVCapability {
	return []db.KVCapability{db.KVMutable}
}
