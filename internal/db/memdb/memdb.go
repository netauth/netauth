// Package memdb is an entirely in-memory database with no persistence
// and no concurrency guarantees.  Its function is to provide a full
// database implementation to tests to run against.
package memdb

import (
	"github.com/hashicorp/go-hclog"

	"github.com/netauth/netauth/internal/db"
	"github.com/netauth/netauth/internal/db/util"
	"github.com/netauth/netauth/internal/health"
	"github.com/netauth/netauth/internal/startup"

	pb "github.com/netauth/protocol"
)

func init() {
	startup.RegisterCallback(cb)
}

func cb() {
	db.Register("MemDB", New)
}

// The MemDB type binds the methods of this "database".  This DB is
// designed really only for supporting the tests of other modules, so
// keep in mind that it is not safe for concurrent execution.
type MemDB struct {
	idx *util.SearchIndex

	eMap map[string]*pb.Entity
	gMap map[string]*pb.Group
}

// New returns a usable memdb with internal structures initialized.
func New(l hclog.Logger) (db.DB, error) {
	l = l.Named("memdb")
	x := &MemDB{
		idx:  util.NewIndex(l),
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
	if ID == "load-error" {
		return nil, db.ErrInternalError
	}
	e, ok := m.eMap[ID]
	if !ok {
		return nil, db.ErrUnknownEntity
	}
	return e, nil
}

// SaveEntity saves an entity to the "database".
func (m *MemDB) SaveEntity(e *pb.Entity) error {
	if e.GetID() == "save-error" {
		return db.ErrInternalError
	}
	m.eMap[e.GetID()] = e
	return m.idx.IndexEntity(e)
}

// DeleteEntity deletes an entity from the "database".
func (m *MemDB) DeleteEntity(ID string) error {
	if _, ok := m.eMap[ID]; !ok {
		return db.ErrUnknownEntity
	}

	delete(m.eMap, ID)
	return m.idx.DeleteEntity(&pb.Entity{ID: &ID})
}

// NextEntityNumber fetches out the next unassigned entity number.
func (m *MemDB) NextEntityNumber() (int32, error) {
	return util.NextEntityNumber(m.LoadEntity, m.DiscoverEntityIDs)
}

// SearchEntities returns a slice of entity given a searchrequest.
func (m *MemDB) SearchEntities(r db.SearchRequest) ([]*pb.Entity, error) {
	res, err := m.idx.SearchEntities(r)
	if err != nil {
		return nil, err
	}
	return util.LoadEntityBatch(res, m.LoadEntity)
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
	if name == "load-error" {
		return nil, db.ErrInternalError
	}
	g, ok := m.gMap[name]
	if !ok {
		return nil, db.ErrUnknownGroup
	}
	return g, nil
}

// SaveGroup saves a group to the "database".
func (m *MemDB) SaveGroup(g *pb.Group) error {
	if g.GetName() == "save-error" {
		return db.ErrInternalError
	}
	m.gMap[g.GetName()] = g
	return m.idx.IndexGroup(g)
}

// DeleteGroup deletes a group from the "database".
func (m *MemDB) DeleteGroup(name string) error {
	if _, ok := m.gMap[name]; !ok {
		return db.ErrUnknownGroup
	}

	delete(m.gMap, name)
	return m.idx.DeleteGroup(&pb.Group{Name: &name})
}

// NextGroupNumber uses the util package to return a group number.
func (m *MemDB) NextGroupNumber() (int32, error) {
	return util.NextGroupNumber(m.LoadGroup, m.DiscoverGroupNames)
}

// SearchGroups returns a slice of entity given a searchrequest.
func (m *MemDB) SearchGroups(r db.SearchRequest) ([]*pb.Group, error) {
	res, err := m.idx.SearchGroups(r)
	if err != nil {
		return nil, err
	}
	return util.LoadGroupBatch(res, m.LoadGroup)
}

func (m *MemDB) healthCheck() health.SubsystemStatus {
	return health.SubsystemStatus{
		OK:     true,
		Name:   "MemDB",
		Status: "MemDB is operating normally",
	}
}
