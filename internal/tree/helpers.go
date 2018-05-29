package tree

import (
	"log"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/internal/crypto"
	"github.com/NetAuth/NetAuth/internal/db"

	pb "github.com/NetAuth/Protocol"
)

// New returns an initialized EMDataStore on to which all other
// functions are bound.
func New(db db.EMDiskInterface, crypto crypto.EMCrypto) *Manager {
	x := Manager{}
	x.bootstrap_done = false
	x.db = db
	x.crypto = crypto
	log.Println("Initialized new Entity Manager")

	return &x
}

// nextUIDNumber computes the next available number to be assigned.
// This allows a NewEntity request to be made with the number field
// unset.
func (m Manager) nextUIDNumber() (int32, error) {
	var largest int32 = 0

	// Iterate over the entities and return the largest ID found
	// +1.  This allows them to be in any order or have IDs
	// missing in the middle and still work.  Though an
	// inefficient search this is worst case O(N) and happends
	// only on provisioning a new entry in the database.
	el, err := m.db.DiscoverEntityIDs()
	if err != nil {
		return 0, err
	}

	for _, en := range el {
		e, err := m.db.LoadEntity(en)
		if err != nil {
			return 0, err
		}
		if e.GetNumber() > largest {
			largest = e.GetNumber()
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

// dedupEntityList takes in a list of entities and deduplicates them
// using a map.
func dedupEntityList(entList []*pb.Entity) []*pb.Entity {
	eMap := make(map[string]*pb.Entity)
	for _, e := range entList {
		eMap[e.GetID()] = e
	}

	// Back to a list...
	var eList []*pb.Entity
	for _, e := range eMap {
		eList = append(eList, e)
	}
	return eList
}

// entityListDifference computes the set of entities that are in list
// a and not in list b.
func entityListDifference(a, b []*pb.Entity) []*pb.Entity {
	diffMap := make(map[string]*pb.Entity)
	// Get a map of the possible options
	for _, e := range a {
		diffMap[e.GetID()] = e
	}
	// Remove the ones that are in the exclude map
	for _, e := range b {
		delete(diffMap, e.GetID())
	}

	var entList []*pb.Entity
	for _, ent := range diffMap {
		entList = append(entList, ent)
	}

	return entList
}
