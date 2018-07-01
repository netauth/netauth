package db

// This package implements a map of interfaces that contain the
// various database options.

import (
	"errors"

	pb "github.com/NetAuth/Protocol"
)

// DB specifies the methods that a DB engine must provide.
type DB interface {
	// Entity handling
	DiscoverEntityIDs() ([]string, error)
	LoadEntity(string) (*pb.Entity, error)
	SaveEntity(*pb.Entity) error
	DeleteEntity(string) error

	// Group handling
	DiscoverGroupNames() ([]string, error)
	LoadGroup(string) (*pb.Group, error)
	SaveGroup(*pb.Group) error
	DeleteGroup(string) error
}

// Factory defines the function which can be used to register new
// implementations.
type Factory func() DB

var (
	backends = make(map[string]Factory)

	// ErrUnknownEntity is returned for requests to load an entity
	// that does not exist.
	ErrUnknownEntity = errors.New("The specified entity does not exist")

	// ErrUnknownGroup is returned for requests to load a group
	// that does not exist.
	ErrUnknownGroup = errors.New("The specified group does not exist")

	// ErrUnknownDatabase is returned for an attempt to create a
	// new database that hasn't been registered.
	ErrUnknownDatabase = errors.New("The specified database does not exist")

	// ErrInternalError is used for all other errors that occur
	// within a database implementation.
	ErrInternalError = errors.New("The database has encountered an internal error")
)

// New returns a db struct.
func New(name string) (DB, error) {
	b, ok := backends[name]
	if !ok {
		return nil, ErrUnknownDatabase
	}
	return b(), nil
}

// RegisterDB takes in a name of the database to register and a
// function signature to bind to that name.
func RegisterDB(name string, newFunc Factory) {
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
