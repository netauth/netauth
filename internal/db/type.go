package db

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
	NextEntityNumber() (int32, error)
	SearchEntities(SearchRequest) ([]*pb.Entity, error)

	// Group handling
	DiscoverGroupNames() ([]string, error)
	LoadGroup(string) (*pb.Group, error)
	SaveGroup(*pb.Group) error
	DeleteGroup(string) error
	NextGroupNumber() (int32, error)
	SearchGroups(SearchRequest) ([]*pb.Group, error)
}

// Factory defines the function which can be used to register new
// implementations.
type Factory func() (DB, error)

// Callback is a function type registered by an external customer that
// is interested in some change that might happen in the storage
// system.  These are returned with a DBEvent populated of whether or
// not the event pertained to an entity or a group, and the relavant
// primary key.
type Callback func(Event)

// Event is a type of message that can be fed to callbacks
// describing the event and the key of the thing that happened.
type Event struct {
	Type string
	PK   string
}

// SearchRequest is an expression that can be interpreted by the
// default util search system, or translated by a storage layer to
// provide a more optimized searching experience.
type SearchRequest struct {
	Expression string
}
