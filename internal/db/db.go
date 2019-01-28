// Package db implements a plugin system for data storage options.
// The db package itself implements the registration and
// initialization functions that provide a uniform interface to
// underlying storage mechanisms.
package db

import (
	"github.com/spf13/viper"
)

var (
	backends map[string]Factory
)

func init() {
	backends = make(map[string]Factory)
}

// New returns a db struct.
func New() (DB, error) {
	b, ok := backends[viper.GetString("db.backend")]
	if !ok {
		return nil, ErrUnknownDatabase
	}
	return b()
}

// Register takes in a name of the database to register and a
// function signature to bind to that name.
func Register(name string, newFunc Factory) {
	if _, ok := backends[name]; ok {
		// Return if the backend is already registered.
		return
	}
	backends[name] = newFunc
}

// GetBackendList returns a string list of the backends that are available
func GetBackendList() []string {
	var l []string

	for b := range backends {
		l = append(l, b)
	}

	return l
}
