package db

// This package implements a map of interfaces that contain the
// various database options.

import (
	"github.com/NetAuth/NetAuth/pkg/errors"

	pb "github.com/NetAuth/NetAuth/pkg/proto"
)

type EMDiskInterface interface {
	// Entity handling
	DiscoverEntityIDs() ([]string, error)
	LoadEntity(string) (*pb.Entity, error)
	LoadEntityNumber(int32) (*pb.Entity, error)
	SaveEntity(*pb.Entity) error
	DeleteEntity(string) error

	// Group handling
	DiscoverGroupNames() ([]string, error)
	LoadGroup(string) (*pb.Group, error)
	LoadGroupNumber(int32) (*pb.Group, error)
	SaveGroup(*pb.Group) error
	DeleteGroup(string) error
}

type DBFactory func() EMDiskInterface

var (
	backends = make(map[string]DBFactory)
)

// NewDB returns a db struct.
func New(name string) (EMDiskInterface, error) {
	b, ok := backends[name]
	if !ok {
		return nil, errors.E_NO_SUCH_DATABASE
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
}

// GetBackendList returns a string list of the backends that are available
func GetBackendList() []string {
	l := []string{}

	for b, _ := range backends {
		l = append(l, b)
	}

	return l
}
