package db

// This package implements a map of interfaces that contain the
// various database options.

import (
	"errors"
	"log"

	"github.com/NetAuth/NetAuth/internal/server/entity_manager"
)

type DBFactory func() entity_manager.EMDiskInterface

var (
	backends           = make(map[string]DBFactory)
	E_NO_SUCH_DATABASE = errors.New("No such database!")
)

// NewDB returns a db struct.
func New(name string) (entity_manager.EMDiskInterface, error) {
	b, ok := backends[name]
	if !ok {
		return nil, E_NO_SUCH_DATABASE
	}
	return b(), nil
}

// RegisterDB takes in a name of the database to register and a
// function signature to bind to that name.
func RegisterDB(name string, newFunc DBFactory) {
	if _, ok := backends[name]; ok {
		// Return if the backend is already registered.
		return
	}
	backends[name] = newFunc
	log.Printf("Registered database implementation '%s'", name)
}

// GetBackendList returns a string list of the backends that are available
func GetBackendList() []string {
	l := []string{}

	for b, _ := range backends {
		l = append(l, b)
	}

	return l
}
