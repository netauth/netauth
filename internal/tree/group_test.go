package tree

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/internal/db/impl/MemDB"
	"github.com/NetAuth/NetAuth/internal/crypto/impl/nocrypto"
	"github.com/NetAuth/NetAuth/pkg/errors"
	pb "github.com/NetAuth/Protocol"
)

func TestNewGroup(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	s := []struct {
		name        string
		displayName string
		gidNumber   int32
		wantErr     error
	}{
		{"fooGroup", "", 1, nil},
		{"fooGroup", "", 1, errors.E_DUPLICATE_GROUP_ID},
		{"barGroup", "", 0, errors.E_DUPLICATE_GROUP_NUMBER},
		{"barGroup", "", 1, errors.E_DUPLICATE_GROUP_NUMBER},
		{"barGroup", "", -1, nil},
	}
	for _, c := range s {
		if err := em.NewGroup(c.name, c.displayName, "", c.gidNumber); err != c.wantErr {
			t.Errorf("Wrong Error: want '%v' got '%v'", c.wantErr, err)
		}
	}
}

func TestDeleteGroup(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	if err := em.NewGroup("foo", "", "", -1); err != nil {
		t.Error(err)
	}

	if _, err := em.GetGroupByName("foo"); err != nil {
		t.Error(err)
	}

	if err := em.DeleteGroup("foo"); err != nil {
		t.Error(err)
	}

	if _, err := em.GetGroupByName("foo"); err != errors.E_NO_GROUP {
		t.Error(err)
	}
}

func TestUpdateGroupMetaInternal(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	if err := em.NewGroup("foo", "foo", "", -1); err != nil {
		t.Error(err)
	}

	update := &pb.Group{DisplayName: proto.String("Foo Group")}

	if err := em.UpdateGroupMeta("foo", update); err != nil {
		t.Error(err)
	}

	g, err := em.GetGroupByName("foo")
	if err != nil {
		t.Error(err)
	}

	if g.GetDisplayName() != "Foo Group" {
		t.Error("Meta update failed!")
	}
}
