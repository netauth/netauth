package entity_manager

import pb "github.com/NetAuth/NetAuth/proto"

type EMDataStore struct {
	eByID map[string]*pb.Entity
	eByUIDNumber map[int32]*pb.Entity

	// Making a bootstrap entity is a rare thing and short
	// circuits most of the permissions logic.  As such we only
	// allow it to be done once per server start.
	bootstrap_done bool

	// The persistence layer contains the functions that actually
	// deal with the disk and make this a useable server.
	db EMDiskInterface
}

type EMDiskInterface interface {
	DiscoverEntityIDs() ([]string, error)
	LoadEntity(string) (*pb.Entity, error)
	SaveEntity(*pb.Entity) error
	DeleteEntity(string) error
}
