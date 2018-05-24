package MemDB

import (
	"github.com/NetAuth/NetAuth/internal/db"
	"github.com/NetAuth/NetAuth/pkg/errors"
	pb "github.com/NetAuth/Protocol"
)

func init() {
	db.RegisterDB("MemDB", New)
}

// This backend is the bare minimum, it just mirrors the
// entity_manager's memory.

type MemDB struct {
	eMap map[string]*pb.Entity
	gMap map[string]*pb.Group
}

func New() db.EMDiskInterface {
	return &MemDB{
		eMap: make(map[string]*pb.Entity),
		gMap: make(map[string]*pb.Group),
	}
}

func (m *MemDB) DiscoverEntityIDs() ([]string, error) {
	var entities []string
	for _, e := range m.eMap {
		entities = append(entities, e.GetID())
	}

	return entities, nil
}

func (m *MemDB) LoadEntity(ID string) (*pb.Entity, error) {
	e, ok := m.eMap[ID]
	if !ok {
		return nil, errors.E_NO_ENTITY
	}
	return e, nil
}

func (m *MemDB) LoadEntityNumber(number int32) (*pb.Entity, error) {
	l, err := m.DiscoverEntityIDs()
	if err != nil {
		return nil, err
	}

	for _, en := range l {
		e, err := m.LoadEntity(en)
		if err != nil {
			return nil, err
		}
		if e.GetUidNumber() == number {
			return e, nil
		}
	}
	return nil, errors.E_NO_ENTITY
}

func (m *MemDB) SaveEntity(e *pb.Entity) error {
	m.eMap[e.GetID()] = e
	return nil
}

func (m *MemDB) DeleteEntity(ID string) error {
	delete(m.eMap, ID)
	return nil
}

func (m *MemDB) DiscoverGroupNames() ([]string, error) {
	var groups []string
	for _, g := range m.gMap {
		groups = append(groups, g.GetName())
	}
	return groups, nil
}

func (m *MemDB) LoadGroup(name string) (*pb.Group, error) {
	g, ok := m.gMap[name]
	if !ok {
		return nil, errors.E_NO_GROUP
	}
	return g, nil
}

func (m *MemDB) LoadGroupNumber(number int32) (*pb.Group, error) {
	l, err := m.DiscoverGroupNames()
	if err != nil {
		return nil, err
	}

	for _, gn := range l {
		g, err := m.LoadGroup(gn)
		if err != nil {
			return nil, err
		}
		if g.GetGidNumber() == number {
			return g, nil
		}
	}
	return nil, errors.E_NO_GROUP
}

func (m *MemDB) SaveGroup(g *pb.Group) error {
	m.gMap[g.GetName()] = g
	return nil
}

func (m *MemDB) DeleteGroup(name string) error {
	if _, ok := m.gMap[name]; !ok {
		return errors.E_NO_GROUP
	}

	delete(m.gMap, name)
	return nil
}
