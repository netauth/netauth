package entity_manager

import pb "github.com/NetAuth/NetAuth/pkg/proto"

type EMDataStore struct {
	// The entities are indexed by ID and uidNumber for
	// convenience in lookup operations that need to be fast.
	eByID        map[string]*pb.Entity
	eByUIDNumber map[int32]*pb.Entity

	// Groups are similarly indexed by ID and gidNumber for
	// convenience in lookup operations that need to be fast.
	gByName      map[string]*pb.Group
	gByGIDNumber map[int32]*pb.Group

	// Making a bootstrap entity is a rare thing and short
	// circuits most of the permissions logic.  As such we only
	// allow it to be done once per server start.
	bootstrap_done bool

	// The persistence layer contains the functions that actually
	// deal with the disk and make this a useable server.
	db EMDiskInterface
}

type EMDiskInterface interface {
	// Entity handling
	DiscoverEntityIDs() ([]string, error)
	LoadEntity(string) (*pb.Entity, error)
	SaveEntity(*pb.Entity) error
	DeleteEntity(string) error

	// Group handling
	// DiscoverGroupIDs() ([]string, error)
	// LoadGroup(string) (*pb.Group, error)
	SaveGroup(*pb.Group) error
	// DeleteGroup(string) error
}
