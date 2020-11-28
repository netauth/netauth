package db

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

// IsEmpty is used to test for an empty event being returned.
func (e *Event) IsEmpty() bool {
	return e.PK == ""
}
