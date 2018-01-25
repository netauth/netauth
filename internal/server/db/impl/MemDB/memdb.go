package MemDB

import (
	"github.com/NetAuth/NetAuth/internal/server/db"
	"github.com/NetAuth/NetAuth/internal/server/entity_manager"
	"github.com/NetAuth/NetAuth/pkg/errors"
	pb "github.com/NetAuth/NetAuth/pkg/proto"
)

func init() {
	db.RegisterDB("MemDB", New)
}

// This is a noop backend that exists almost exclusively to prove that
// the backends work.  This backend noops all the functions since it
// doesn't add value to store things in memory twice (i.e. the cache
// in the entity_manager itself is now treated as the primary database.

type MemDB struct{}

func New() entity_manager.EMDiskInterface                { return &MemDB{} }
func (m *MemDB) DiscoverEntityIDs() ([]string, error)    { return []string{}, nil }
func (m *MemDB) LoadEntity(_ string) (*pb.Entity, error) { return nil, errors.E_NO_ENTITY }
func (m *MemDB) SaveEntity(_ *pb.Entity) error           { return nil }
func (m *MemDB) DeleteEntity(_ string) error             { return nil }
