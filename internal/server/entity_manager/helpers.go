package entity_manager

import (
	"log"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/internal/server/health"
	"github.com/NetAuth/NetAuth/pkg/errors"
	pb "github.com/NetAuth/NetAuth/pkg/proto"
)

// New returns an initialized EMDataStore on to which all other
// functions are bound.
func New(db *EMDiskInterface) *EMDataStore {
	x := EMDataStore{}
	x.bootstrap_done = false
	if db != nil {
		x.db = *db
	}
	log.Println("Initialized new Entity Manager")

	if x.db == nil {
		log.Println("Entity Manager persistence layer is not available!")
	}

	return &x
}

// initMem sets up the memory of the datastore
func (emds *EMDataStore) initMem() {
	emds.eByID = make(map[string]*pb.Entity)
	emds.eByUIDNumber = make(map[int32]*pb.Entity)
	emds.gByName = make(map[string]*pb.Group)
	emds.gByGIDNumber = make(map[int32]*pb.Group)
}

// Reload conducts an in place swap of the entity_manager and causes
// it to reconcile state with what's in long term storage.
func (emds *EMDataStore) Reload() {
	// If we don't have any kind of backing database, then this
	// should be noop'd out.
	if emds.db == nil {
		return
	}

	log.Println("Beginning EM Reload")

	// We're about to dump the caches, so lets make sure we mark
	// ourselves bad first.
	health.SetBad()

	// Reset the internal memory
	emds.initMem()

	// Get the entity list from disk
	el, err := emds.db.DiscoverEntityIDs()
	if err != nil {
		log.Printf("Cannot reload the entity_manager (%s)", err)
		return
	}

	// Load entities in from the disk.
	loaded := 0
	for _, en := range el {
		emds.loadFromDisk(en)
		loaded++
	}

	// Log the reload data
	log.Printf("Discovered %d entities; loaded %d", len(el), loaded)
	if len(el) != loaded {
		// If we have the wrong number, then return now
		// without marking the server healthy again.
		return
	}

	// Everythings reloaded, we'll mark healthy again
	log.Println("EM Reload Complete!")
	health.SetGood()
}

// nextUIDNumber computes the next available uidNumber to be assigned.
// This allows a NewEntity request to be made with the uidNumber field
// unset.
func (emds *EMDataStore) nextUIDNumber() int32 {
	var largest int32 = 0

	// Iterate over the entities and return the largest ID found
	// +1.  This allows them to be in any order or have IDs
	// missing in the middle and still work.  Though an
	// inefficient search this is worst case O(N) and happens in
	// memory anyway.
	for i := range emds.eByID {
		if emds.eByID[i].GetUidNumber() > largest {
			largest = emds.eByID[i].GetUidNumber()
		}
	}

	return largest + 1
}

// getEntityByID returns a pointer to an Entity struct and an error
// value.  The error value will either be errors.E_NO_ENTITY if the
// requested value did not match, or will be nil where an entity is
// returned.  The string must be a complete match for the entity name
// being requested.
func (emds *EMDataStore) getEntityByID(ID string) (*pb.Entity, error) {
	e, ok := emds.eByID[ID]
	if !ok {
		// Attempt to load the entity if the persistence layer
		// is available.
		if emds.db != nil {
			if err := emds.loadFromDisk(ID); err != nil {
				return nil, err
			}
		}
	}

	// Try again after potentially having loaded the entity
	e, ok = emds.eByID[ID]
	if !ok {
		return nil, errors.E_NO_ENTITY
	}
	return e, nil
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
