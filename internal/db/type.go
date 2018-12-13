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

// SearchRequest is an expression that can be interpreted by the
// default util search system, or translated by a storage layer to
// provide a more optimized searching experience.
type SearchRequest struct {
	Expression string
}
