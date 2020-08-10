// Package db implements a plugin system for data storage options.
// The db package itself implements the registration and
// initialization functions that provide a uniform interface to
// underlying storage mechanisms.
package db

import (
	"github.com/hashicorp/go-hclog"
)

var (
	lb       hclog.Logger
	backends map[string]Factory
)

func init() {
	backends = make(map[string]Factory)
}

// New returns a db struct.
func New(backend string) (DB, error) {
	b, ok := backends[backend]
	if !ok {
		return nil, ErrUnknownDatabase
	}
	log().Info("Initializing database backend", "backend", backend)
	return b(log())
}

// Register takes in a name of the database to register and a
// function signature to bind to that name.
func Register(name string, newFunc Factory) {
	if _, ok := backends[name]; ok {
		// Return if the backend is already registered.
		return
	}
	backends[name] = newFunc
	log().Debug("Registered backend", "backend", name)
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
