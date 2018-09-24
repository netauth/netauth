package db

// This package implements a map of interfaces that contain the
// various database options.

import (
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
type Factory func() (DB, error)

var (
	backends map[string]Factory
)

func init() {
	backends = make(map[string]Factory)
}

// New returns a db struct.
func New(name string) (DB, error) {
	b, ok := backends[name]
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
