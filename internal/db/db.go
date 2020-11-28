// Package db implements a uniform mechanism for interacting with
// entities and groups on top of a generic key/value store which is
// used for persistent data.
package db

import (
	"path"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-hclog"

	types "github.com/netauth/protocol"
)

var (
	lb hclog.Logger
)

// New returns a db struct.
func New(backend string) (*DB, error) {
	kv, err := newKV(backend, log())
	if err != nil {
		return nil, err
	}

	idx := NewIndex(log())
	x := &DB{
		log:   log(),
		Index: idx,
		kv:    kv,
		cbs:   make(map[string]Callback),
	}
	kv.SetEventFunc(x.FireEvent)
	x.Index.ConfigureCallback(x.LoadEntity, x.LoadGroup)
	x.RegisterCallback("BleveSearch", x.Index.IndexCallback)
	return x, nil
}

// DiscoverEntityIDs searches the keyspace for all entity IDs.  All
// returned strings are loadable entities.
func (db *DB) DiscoverEntityIDs() ([]string, error) {
	return db.kv.Keys("/entities/*")
}

// LoadEntity retrieves a single entity from the kv store.
func (db *DB) LoadEntity(ID string) (*types.Entity, error) {
	b, err := db.kv.Get(path.Join("/entities", ID))
	if err == ErrNoValue {
		return nil, ErrUnknownEntity
	}
	if err != nil {
		db.log.Debug("Error loading entity from KV store", "error", err, "ID", ID)
		return nil, err
	}

	e := &types.Entity{}
	if err := proto.Unmarshal(b, e); err != nil {
		db.log.Warn("Error unmarshaling entity", "error", err)
		return nil, err
	}
	return e, nil
}

// SaveEntity writes an entity to the kv store.
func (db *DB) SaveEntity(e *types.Entity) error {
	b, err := proto.Marshal(e)
	if err != nil {
		db.log.Warn("Error marshaling entity", "error", err)
		return err
	}

	if err := db.kv.Put(path.Join("/entities", e.GetID()), b); err != nil {
		db.log.Warn("Error storing entity", "error", err)
		return err
	}
	return nil
}

// DeleteEntity tries to delete an entity that already exists.
func (db *DB) DeleteEntity(ID string) error {
	err := db.kv.Del(path.Join("/entities", ID))
	if err == ErrNoValue {
		return ErrUnknownEntity
	}
	return err
}

// DiscoverGroupNames searches the keyspace for all group names.  All
// returned strings are loadable groups.
func (db *DB) DiscoverGroupNames() ([]string, error) {
	return db.kv.Keys("/groups/*")
}

// LoadGroup retrieves a single group from the kv store.
func (db *DB) LoadGroup(ID string) (*types.Group, error) {
	b, err := db.kv.Get(path.Join("/groups", ID))
	if err == ErrNoValue {
		return nil, ErrUnknownGroup
	}
	if err != nil {
		db.log.Debug("Error loading group from KV store", "error", err, "ID", ID)
		return nil, err
	}

	g := &types.Group{}
	if err := proto.Unmarshal(b, g); err != nil {
		db.log.Warn("Error unmarshaling group", "error", err)
		return nil, err
	}
	return g, nil
}

// SaveGroup writes an group to the kv store.
func (db *DB) SaveGroup(g *types.Group) error {
	b, err := proto.Marshal(g)
	if err != nil {
		db.log.Warn("Error marshaling group", "error", err)
		return err
	}

	if err := db.kv.Put(path.Join("/groups", g.GetName()), b); err != nil {
		db.log.Warn("Error storing group", "error", err)
		return err
	}
	return nil
}

// DeleteGroup tries to delete an group that already exists.
func (db *DB) DeleteGroup(ID string) error {
	err := db.kv.Del(path.Join("/groups", ID))
	if err == ErrNoValue {
		return ErrUnknownGroup
	}
	return err
}

// Shutdown is called to disconnect the KV store from any other
// systems and flush any buffers before shutting down the server.
func (db *DB) Shutdown() {
	if err := db.kv.Close(); err != nil {
		db.log.Error("Error shutting down KV store", "error", err)
	}
}

// NextEntityNumber computes and returns the next unnassigned number
// in the entity space.
func (db *DB) NextEntityNumber() (int32, error) {
	var largest int32

	// Iterate over the entities and return the largest ID found
	// +1.  This allows them to be in any order or have IDs
	// missing in the middle and still work.  Though an
	// inefficient search this is worst case O(N) and happens only
	// on provisioning a new entry in the database.
	el, err := db.DiscoverEntityIDs()
	if err != nil {
		return 0, err
	}

	for _, en := range el {
		e, err := db.LoadEntity(path.Base(en))
		if err != nil {
			return 0, err
		}
		if e.GetNumber() > largest {
			largest = e.GetNumber()
		}
	}

	return largest + 1, nil
}

// NextGroupNumber computes the next available group number and
// returns it.
func (db *DB) NextGroupNumber() (int32, error) {
	var largest int32

	l, err := db.DiscoverGroupNames()
	if err != nil {
		return 0, err
	}
	for _, i := range l {
		g, err := db.LoadGroup(path.Base(i))
		if err != nil {
			return 0, err
		}
		if g.GetNumber() > largest {
			largest = g.GetNumber()
		}
	}

	return largest + 1, nil
}

// Capabilities returns a slice of capabilities the backing store
// supports.  This allows higher level abstractions to decide if they
// want to return errors in certain circumstances, such as this
// instance not being writeable.
func (db *DB) Capabilities() []KVCapability {
	return db.kv.Capabilities()
}

// SearchEntities performs a search of all entities using the given
// query and then batch loads the result.
func (db *DB) SearchEntities(r SearchRequest) ([]*types.Entity, error) {
	ids, err := db.Index.SearchEntities(r)
	if err != nil {
		return nil, err
	}

	return db.loadEntityBatch(ids)
}

// SearchGroups performs a search of all groups using the given query
// and then batch loads the result.
func (db *DB) SearchGroups(r SearchRequest) ([]*types.Group, error) {
	ids, err := db.Index.SearchGroups(r)
	if err != nil {
		return nil, err
	}

	return db.loadGroupBatch(ids)
}

func (db *DB) loadEntityBatch(ids []string) ([]*types.Entity, error) {
	eSlice := []*types.Entity{}

	for i := range ids {
		e, err := db.LoadEntity(ids[i])
		if err != nil {
			return nil, err
		}
		eSlice = append(eSlice, e)
	}
	return eSlice, nil
}

func (db *DB) loadGroupBatch(ids []string) ([]*types.Group, error) {
	gSlice := []*types.Group{}

	for i := range ids {
		g, err := db.LoadGroup(ids[i])
		if err != nil {
			return nil, err
		}
		gSlice = append(gSlice, g)
	}
	return gSlice, nil
}

// SetParentLogger sets the parent logger for this instance.
func SetParentLogger(l hclog.Logger) {
	lb = l.Named("db")
}

// log is a convenience function that will return a null logger if a
// parent logger has not been specified, mostly useful for tests.
func log() hclog.Logger {
	if lb == nil {
		lb = hclog.NewNullLogger()
	}
	return lb
}
