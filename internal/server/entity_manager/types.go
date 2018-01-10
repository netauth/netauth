package entity_manager

import pb "github.com/NetAuth/NetAuth/proto"

type EMDataStore struct {
	eByID map[string]*pb.Entity
	eByUIDNumber map[int32]*pb.Entity

	// Making a bootstrap entity is a rare thing and short
	// circuits most of the permissions logic.  As such we only
	// allow it to be done once per server start.
	bootstrap_done bool
}
