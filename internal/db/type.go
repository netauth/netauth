package db

import (
	"github.com/hashicorp/go-hclog"

	types "github.com/netauth/protocol"
)

// KVFactory returns a KVStore, and is a registeryable function during
// init to be called later.
type KVFactory func(hclog.Logger) (KVStore, error)

// A KVStore is the backing mechanism that deals with persisting data
// to somewhere that won't lose it.  This can be the disk, a remote
// blob store, the desk of a particularly trusted employee, etc.
type KVStore interface {
	Put(string, []byte) error
	Get(string) ([]byte, error)
	Del(string) error

	Keys(string) ([]string, error)
	Close() error

	Capabilities() []KVCapability
	SetEventFunc(func(Event))
}

// A DB is a collection of methods satisfying tree.DB, and which read
// and write data to a KVStore
type DB struct {
	log hclog.Logger
	kv  KVStore
	cbs map[string]Callback

	*Index
}

// A KVCapability is a specific property that a KV Store might have.
// It allows stores to express things like supporting HA access.
type KVCapability int

// Callback is a function type registered by an external customer that
// is interested in some change that might happen in the storage
// system.  These are returned with a DBEvent populated of whether or
// not the event pertained to an entity or a group, and the relavant
// primary key.
type Callback func(Event)

// Event is a type of message that can be fed to callbacks
// describing the event and the key of the thing that happened.
type Event struct {
	Type EventType
	PK   string
}

// An EventType is used to specify what kind of event has happened and
// is constrained for consumption in downstream select cases.  As
// these IDs are entirely internal and are maintained within a
// compiled version, iota is used here to make it easier to patch this
// list in the future.
type EventType int

// The callbacks defined below are used to signal what events are
// handled by the Event subsystem.
const (
	EventEntityCreate EventType = iota
	EventEntityUpdate
	EventEntityDestroy

	EventGroupCreate
	EventGroupUpdate
	EventGroupDestroy
)

// SearchRequest is an expression that can be interpreted by the
// default util search system, or translated by a storage layer to
// provide a more optimized searching experience.
type SearchRequest struct {
	Expression string
}

// These allow the index to get limited access to the db itself.  You
// might ask why we can't embed an interface here, and its because the
// database embeds an index, so this would create an embed cycle which
// is not allowed.
type loadEntityFunc func(string) (*types.Entity, error)
type loadGroupFunc func(string) (*types.Group, error)
