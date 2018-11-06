package tree

import (
	"github.com/NetAuth/NetAuth/internal/crypto"
	"github.com/NetAuth/NetAuth/internal/db"

	pb "github.com/NetAuth/Protocol"
)

// The Manager binds all methods for managing a tree of entities with
// the associated groups, capabilities, and other assorted functions.
// This is the type that is served up by the RPC layer.
type Manager struct {
	// Making a bootstrap entity is a rare thing and short
	// circuits most of the permissions logic.  As such we only
	// allow it to be done once per server start.
	bootstrapDone bool

	// The persistence layer contains the functions that actually
	// deal with the disk and make this a useable server.
	db db.DB

	// The Crypto layer allows us to plug in different crypto
	// engines
	crypto crypto.EMCrypto

	processors map[string]EntityProcessor

	entityProcesses map[string][]EntityProcessorHook
}

// An EntityProcessor is a chain of functions that modify entities in
// some way.
type EntityProcessor struct {
	Entity      *pb.Entity
	RequestData *pb.Entity
	hooks       []EntityProcessorHook
}

// An EntityProcessorHook is a function that transforms an entity as
// part of an EntityProcessor pipeline.
type EntityProcessorHook interface {
	Priority() int
	Name() string
	Run(*pb.Entity, *pb.Entity) error
}
