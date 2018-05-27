package tree

import (
	"log"
	"sort"

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

// nextUIDNumber computes the next available uidNumber to be assigned.
// This allows a NewEntity request to be made with the uidNumber field
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
	// Sort the slices, this makes the end check a bit faster
	sort.Slice(a, func(i, j int) bool {
		return a[i].GetUidNumber() < a[j].GetUidNumber()
	})
	sort.Slice(b, func(i, j int) bool {
		return b[i].GetUidNumber() < b[j].GetUidNumber()
	})

	// Iterate over both lists and pick the ones that are in A and
	// not in B.
	var entList []*pb.Entity
	for i, j := 0, 0; i < len(a); i++ {
		if a[i].GetUidNumber() != b[j].GetUidNumber() && j+1 < len(b) && a[i].GetUidNumber() < b[j+1].GetUidNumber() {
			entList = append(entList, a[i])
		}
		if j+1 < len(b) {
			j++
		}
		if j == len(b)-1 && i < len(a)-1 {
			// we're at the end of list B, so if i is less
			// than len(a) we take from there to the end
			entList = append(entList, a[i:len(a)-1]...)
		}
	}

	return entList
}
