package interface_test

import (
	"context"
	"testing"

	"google.golang.org/protobuf/proto"

	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

// buildSampleTree builds an initial tree that has certain properties
// that are useful for testing in.  Some assertions on this tree are
// stable and allow us to test membership and nesting properties more
// easily, and this tree is also a generally stable thing to test
// against to make sure that the system is reporting sane values.
//
// The tree consists of 3 entities and 5 groups.  The groups are
// arranged as follows:
//
// * group1 has no rules, it directly has members entity1 and entity2
//
// * group2 has an include of group3, it has members entity1 and
// transitively entity3 via group3
//
// * group3 has entity3 as a direct member, no expansion rules are present
//
// * group4 has an include rule for group1 and an exclude rule for
// group5.  It should have logical members entity1 and NOT entity2,
// which is a member of group5
//
// * group5 has no rules and a direct member of entity2
func buildSampleTree(t *testing.T, db tree.DB) {
	ctxt := context.Background()
	group1 := &pb.Group{
		Name: proto.String("group1"),
	}
	group2 := &pb.Group{
		Name: proto.String("group2"),
		Expansions: []string{
			"INCLUDE:group3",
		},
	}
	group3 := &pb.Group{
		Name: proto.String("group3"),
	}
	group4 := &pb.Group{
		Name: proto.String("group4"),
		Expansions: []string{
			"INCLUDE:group1",
			"EXCLUDE:group5",
		},
	}
	group5 := &pb.Group{
		Name: proto.String("group5"),
	}

	if err := db.SaveGroup(ctxt, group1); err != nil {
		t.Fatal(err)
	}
	if err := db.SaveGroup(ctxt, group2); err != nil {
		t.Fatal(err)
	}
	if err := db.SaveGroup(ctxt, group3); err != nil {
		t.Fatal(err)
	}
	if err := db.SaveGroup(ctxt, group4); err != nil {
		t.Fatal(err)
	}
	if err := db.SaveGroup(ctxt, group5); err != nil {
		t.Fatal(err)
	}

	entity1 := &pb.Entity{
		ID: proto.String("entity1"),
		Meta: &pb.EntityMeta{
			Groups: []string{
				"group1",
				"group2",
			},
		},
	}
	entity2 := &pb.Entity{
		ID: proto.String("entity2"),
		Meta: &pb.EntityMeta{
			Groups: []string{
				"group1",
				"group5",
			},
		},
	}
	entity3 := &pb.Entity{
		ID: proto.String("entity3"),
		Meta: &pb.EntityMeta{
			Groups: []string{
				"group3",
			},
		},
	}
	if err := db.SaveEntity(ctxt, entity1); err != nil {
		t.Fatal(err)
	}
	if err := db.SaveEntity(ctxt, entity2); err != nil {
		t.Fatal(err)
	}
	if err := db.SaveEntity(ctxt, entity3); err != nil {
		t.Fatal(err)
	}
}
