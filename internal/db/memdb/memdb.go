package memdb

import (
	"github.com/NetAuth/NetAuth/internal/db"
	"github.com/NetAuth/NetAuth/internal/health"

	pb "github.com/NetAuth/Protocol"
)

func init() {
	db.RegisterDB("MemDB", New)
}

// The MemDB type binds the methods of this "database".  This DB is
// designed really only for supporting the tests of other modules, so
// keep in mind that it is not safe for concurrent execution.
type MemDB struct {
	eMap map[string]*pb.Entity
	gMap map[string]*pb.Group
}

// New returns a usable memdb with internal structures initialized.
func New() (db.DB, error) {
	x := &MemDB{
		eMap: make(map[string]*pb.Entity),
		gMap: make(map[string]*pb.Group),
	}

	health.RegisterCheck("MemDB", x.healthCheck)
	return x, nil
}

// DiscoverEntityIDs returns a list of entity IDs which can then be
// used to load particular entities.
func (m *MemDB) DiscoverEntityIDs() ([]string, error) {
	var entities []string
	for _, e := range m.eMap {
		entities = append(entities, e.GetID())
	}

	return entities, nil
}

// LoadEntity loads an entity from the "database".
func (m *MemDB) LoadEntity(ID string) (*pb.Entity, error) {
	e, ok := m.eMap[ID]
	if !ok {
		return nil, db.ErrUnknownEntity
	}
	return e, nil
}

// SaveEntity saves an entity to the "database".
func (m *MemDB) SaveEntity(e *pb.Entity) error {
	m.eMap[e.GetID()] = e
	return nil
}

// DeleteEntity deletes an entity from the "database".
func (m *MemDB) DeleteEntity(ID string) error {
	if _, ok := m.eMap[ID]; !ok {
		return db.ErrUnknownEntity
	}

	delete(m.eMap, ID)
	return nil
}

// DiscoverGroupNames returns  a slice  of strings  that can  be later
// used to load groups.
func (m *MemDB) DiscoverGroupNames() ([]string, error) {
	var groups []string
	for _, g := range m.gMap {
		groups = append(groups, g.GetName())
	}
	return groups, nil
}

// LoadGroup loads a group from the "database".
func (m *MemDB) LoadGroup(name string) (*pb.Group, error) {
	g, ok := m.gMap[name]
	if !ok {
		return nil, db.ErrUnknownGroup
	}
	return g, nil
}

// SaveGroup saves a group to the "database".
func (m *MemDB) SaveGroup(g *pb.Group) error {
	m.gMap[g.GetName()] = g
	return nil
}

// DeleteGroup deletes a group from the "database".
func (m *MemDB) DeleteGroup(name string) error {
	if _, ok := m.gMap[name]; !ok {
		return db.ErrUnknownGroup
	}

	delete(m.gMap, name)
	return nil
}


func (m *MemDB) healthCheck() health.SubsystemStatus {
	return health.SubsystemStatus{
		OK: true,
		Name: "MemDB",
		Status: "MemDB is operating normally",
	}
}
