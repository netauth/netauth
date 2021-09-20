package db

import (
	"context"
	"path"
)

// RegisterCallback takes a callback name and handle and registers
// them for later calling.
func (db *DB) RegisterCallback(name string, c Callback) {
	if _, ok := db.cbs[name]; ok {
		// Already here...
		log().Warn("Attempted to register duplicate callback", "callback", name)
		return
	}
	db.cbs[name] = c
	log().Info("Database callback registered", "callback", name)
}

// FireEvent fires an event to all callbacks.
func (db *DB) FireEvent(e Event) {
	log().Debug("Processing callbacks")
	for name, c := range db.cbs {
		log().Trace("Calling callback", "callback", name)
		c(e)
	}
}

// EventUpdateAll fires an event for all entities and all groups with
// the type set to "Update".  This is used to allow the async
// components that are event driven to pre-load on a server startup
// and begin monitoring changes after the load completes.
func (db *DB) EventUpdateAll() error {
	ids, err := db.DiscoverEntityIDs(context.Background())
	if err != nil {
		return err
	}
	for _, i := range ids {
		db.FireEvent(Event{Type: EventEntityUpdate, PK: path.Base(i)})
	}

	ids, err = db.DiscoverGroupNames(context.Background())
	if err != nil {
		return err
	}
	for _, i := range ids {
		db.FireEvent(Event{Type: EventGroupUpdate, PK: path.Base(i)})
	}
	return nil
}

// IsEmpty is used to test for an empty event being returned.
func (e *Event) IsEmpty() bool {
	return e.PK == ""
}
