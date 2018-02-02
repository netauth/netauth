package entity_manager

import (
	"log"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/internal/server/db"
	"github.com/NetAuth/NetAuth/internal/server/crypto"

	pb "github.com/NetAuth/NetAuth/pkg/proto"
)

// New returns an initialized EMDataStore on to which all other
// functions are bound.
func New(db db.EMDiskInterface, crypto crypto.EMCrypto) *EMDataStore {
	x := EMDataStore{}
	x.bootstrap_done = false
	x.db = db
	x.crypto = crypto
	log.Println("Initialized new Entity Manager")

	return &x
}

// nextUIDNumber computes the next available uidNumber to be assigned.
// This allows a NewEntity request to be made with the uidNumber field
// unset.
func (emds *EMDataStore) nextUIDNumber() (int32, error) {
	var largest int32 = 0

	// Iterate over the entities and return the largest ID found
	// +1.  This allows them to be in any order or have IDs
	// missing in the middle and still work.  Though an
	// inefficient search this is worst case O(N) and happends
	// only on provisioning a new entry in the database.
	el, err := emds.db.DiscoverEntityIDs()
	if err != nil {
		return 0, err
	}

	for _, en := range el {
		e, err := emds.db.LoadEntity(en)
		if err != nil {
			return 0, err
		}
		if e.GetUidNumber() > largest {
			largest = e.GetUidNumber()
		}
	}

	return largest + 1, nil
}

// safeCopyEntity makes a copy of the entity provided but removes
// fields that are related to security.  This permits the entity that
// is returned to be handed off outside the server.
func safeCopyEntity(e *pb.Entity) (*pb.Entity, error) {
	// Marshal the proto to get a pure data representation of it.
	data, err := proto.Marshal(e)
	if err != nil {
		return nil, err
	}

	// Unmarshaling here ensures that the new entity has no
	// connection to the old one.
	ne := &pb.Entity{}
	if err := proto.Unmarshal(data, ne); err != nil {
		return nil, err
	}

	// Before returning, fields related to security are nulled out
	// so that they aren't available in the returned copy.  At
	// least not available in a meaningful sense.
	ne.Secret = proto.String("<REDACTED>")

	return ne, nil
}
