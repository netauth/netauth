package db

import (
	"testing"

	pb "github.com/NetAuth/Protocol"
)

type dummyDB struct{}

func (*dummyDB) DiscoverEntityIDs() ([]string, error)  { return []string{}, nil }
func (*dummyDB) LoadEntity(string) (*pb.Entity, error) { return nil, nil }
func (*dummyDB) SaveEntity(*pb.Entity) error           { return nil }
func (*dummyDB) DeleteEntity(string) error             { return nil }
func (*dummyDB) DiscoverGroupNames() ([]string, error) { return []string{}, nil }
func (*dummyDB) LoadGroup(string) (*pb.Group, error)   { return nil, nil }
func (*dummyDB) SaveGroup(*pb.Group) error             { return nil }
func (*dummyDB) DeleteGroup(string) error              { return nil }
func newDummyDB() (DB, error)                          { return new(dummyDB), nil }

func TestRegisterDB(t *testing.T) {
	backends = make(map[string]Factory)

	Register("dummy", newDummyDB)
	if l := GetBackendList(); len(l) != 1 || l[0] != "dummy" {
		t.Error("Database factory failed to register")
	}

	Register("dummy", newDummyDB)
	if l := GetBackendList(); len(l) != 1 {
		t.Error("A duplicate database was registered")
	}
}

func TestNewKnown(t *testing.T) {
	backends = make(map[string]Factory)

	Register("dummy", newDummyDB)

	x, err := New("dummy")
	if err != nil {
		t.Error(err)
	}

	if _, ok := x.(*dummyDB); !ok {
		t.Error("Something that isn't a database came out...")
	}
}

func TestNewUnknown(t *testing.T) {
	backends = make(map[string]Factory)
	x, err := New("unknown")
	if x != nil && err != ErrUnknownDatabase {
		t.Error(err)
	}
}
