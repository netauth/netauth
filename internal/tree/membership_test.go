package tree

import (
	"testing"

	"github.com/NetAuth/NetAuth/internal/db"

	pb "github.com/NetAuth/Protocol"
)

func TestMembershipEdit(t *testing.T) {
	em := getNewEntityManager(t)

	e := &pb.Entity{}

	if err := em.NewGroup("fooGroup", "", "", 1000); err != nil {
		t.Error(err)
	}

	if err := em.addEntityToGroup(e, "fooGroup"); err != nil {
		t.Error(err)
	}

	groups := em.getDirectGroups(e)
	if len(groups) != 1 || groups[0] != "fooGroup" {
		t.Error("Wrong group number/membership")
	}

	em.removeEntityFromGroup(e, "fooGroup")
	groups = em.getDirectGroups(e)
	if len(groups) != 0 {
		t.Error("Wrong group number/membership")
	}
}

func TestAddEntityToGroupExternal(t *testing.T) {
	em := getNewEntityManager(t)

	s := []struct {
		entity  string
		group   string
		wantErr error
	}{
		{"foo", "", db.ErrUnknownGroup},
		{"", "", db.ErrUnknownEntity},
		{"foo", "bar", nil},
	}

	if err := em.NewEntity("foo", -1, ""); err != nil {
		t.Fatal(err)
	}
	if err := em.NewGroup("bar", "", "", -1); err != nil {
		t.Fatal(err)
	}

	for i, c := range s {
		if err := em.AddEntityToGroup(c.entity, c.group); err != c.wantErr {
			t.Fatalf("Test %d: Got %v Want %v", i, err, c.wantErr)
		}
	}
}

func TestAddEntityToGroupTwice(t *testing.T) {
	em := getNewEntityManager(t)

	if err := em.NewEntity("foo", -1, ""); err != nil {
		t.Fatal(err)
	}
	if err := em.NewGroup("bar", "", "", -1); err != nil {
		t.Fatal(err)
	}

	e, err := em.GetEntity("foo")
	if err != nil {
		t.Fatal(err)
	}

	// Try a duplicate add and make sure it doesn't explode
	if err := em.addEntityToGroup(e, "bar"); err != nil {
		t.Fatal(err)
	}
	if err := em.addEntityToGroup(e, "bar"); err != nil {
		t.Fatal(err)
	}
}

func TestRemoveEntityFromGroupExternalNoEntity(t *testing.T) {
	em := getNewEntityManager(t)

	if err := em.RemoveEntityFromGroup("", ""); err != db.ErrUnknownEntity {
		t.Fatal(err)
	}
}

func TestRemoveEntityFromGroupExternal(t *testing.T) {
	em := getNewEntityManager(t)

	if err := em.NewEntity("foo", -1, ""); err != nil {
		t.Fatal(err)
	}

	if err := em.RemoveEntityFromGroup("foo", ""); err != nil {
		t.Fatal(err)
	}
}

func TestRemoveEntityFromGroup(t *testing.T) {
	em := getNewEntityManager(t)

	// Get an entity and some groups
	if err := em.NewEntity("foo", -1, ""); err != nil {
		t.Fatal(err)
	}
	if err := em.NewGroup("grp1", "", "", -1); err != nil {
		t.Fatal(err)
	}
	if err := em.NewGroup("grp2", "", "", -1); err != nil {
		t.Fatal(err)
	}
	e, err := em.GetEntity("foo")
	if err != nil {
		t.Fatal(err)
	}

	// Put the entity in some groups
	if err := em.addEntityToGroup(e, "grp1"); err != nil {
		t.Fatal(err)
	}
	if err := em.addEntityToGroup(e, "grp2"); err != nil {
		t.Fatal(err)
	}

	// Pull the entity out of grp1
	if err := em.removeEntityFromGroup(e, "grp1"); err != nil {
		t.Fatal(err)
	}

	// Verify continued membership in grp2
	if len(e.GetMeta().GetGroups()) != 1 || e.GetMeta().GetGroups()[0] != "grp2" {
		t.Fatal("Group membership is wrong")
	}
}

func TestRemoveEntityFromGroupNilMeta(t *testing.T) {
	em := getNewEntityManager(t)

	e := &pb.Entity{}

	// This is just to make sure that this function doesn't
	// explode.
	em.removeEntityFromGroup(e, "fooGroup")
}

func TestGetGroupsNoMeta(t *testing.T) {
	em := getNewEntityManager(t)

	e := &pb.Entity{}

	if groups := em.getDirectGroups(e); len(groups) != 0 {
		t.Error("getDirectGroups fabricated a group!")
	}
}

func TestGetMemberships(t *testing.T) {
	em := getNewEntityManager(t)

	// Create some groups
	if err := em.NewGroup("grp1", "", "", -1); err != nil {
		t.Fatal(err)
	}
	if err := em.NewGroup("grp2", "", "", -1); err != nil {
		t.Fatal(err)
	}

	// Expand grp1 to include grp2
	if err := em.ModifyGroupExpansions("grp1", "grp2", pb.ExpansionMode_INCLUDE); err != nil {
		t.Fatal(err)
	}

	// Create an entity and make them a member of grp2
	if err := em.NewEntity("foo", -1, ""); err != nil {
		t.Fatal(err)
	}
	if err := em.AddEntityToGroup("foo", "grp2"); err != nil {
		t.Fatal(err)
	}

	// Get the memberships of foo
	e, err := em.GetEntity("foo")
	if err != nil {
		t.Fatal(err)
	}
	memberships := em.GetMemberships(e, true)

	// Memberships should be an array of 2 with one of them being
	// grp1 and one of them being grp2
	if len(memberships) != 2 {
		t.Fatal("Wrong number of groups after expansion")
	}

	if memberships[0] != "grp1" && memberships[0] != "grp2" {
		t.Fatal("Wrong group expansions")
	}
	if memberships[1] != "grp1" && memberships[1] != "grp2" {
		t.Fatal("Wrong group expansions")
	}

	// Now we check the same thing but without running the
	// expansions
	memberships = em.GetMemberships(e, false)
	if len(memberships) != 1 && memberships[0] != "grp2" {
		t.Fatal("Wrong number of groups after expansion")
	}
}

func TestListMembersALLInternal(t *testing.T) {
	em := getNewEntityManager(t)

	s := []string{
		"foo",
		"bar",
		"baz",
	}

	for _, id := range s {
		if err := em.NewEntity(id, -1, ""); err != nil {
			t.Error(err)
		}

		listALL, err := em.listMembers("ALL")
		if err != nil {
			t.Error(err)
		}

		dbAll, err := em.db.DiscoverEntityIDs()
		if err != nil {
			t.Error(err)
		}
		if len(dbAll) != len(listALL) {
			t.Error("Different number of entities returned!")
		}
	}
}

func TestListMembersNoMatchInternal(t *testing.T) {
	em := getNewEntityManager(t)
	list, err := em.listMembers("")
	if list != nil && err != db.ErrUnknownGroup {
		t.Error(err)
	}
}

func TestListMembersNoMatchExternal(t *testing.T) {
	em := getNewEntityManager(t)
	list, err := em.ListMembers("")
	if list != nil && err != db.ErrUnknownGroup {
		t.Error(err)
	}
}

func TestListMembersNoExpansions(t *testing.T) {
	em := getNewEntityManager(t)

	if err := em.NewGroup("grp1", "", "", -1); err != nil {
		t.Fatal(err)
	}
	if err := em.NewEntity("foo", -1, ""); err != nil {
		t.Fatal(err)
	}

	if err := em.AddEntityToGroup("foo", "grp1"); err != nil {
		t.Fatal(err)
	}

	members, err := em.listMembers("grp1")
	if err != nil {
		t.Fatal(err)
	}

	if len(members) != 1 || members[0].GetID() != "foo" {
		t.Fatal("Unexpected Membership")
	}
}

func TestListMembersWithExpansions(t *testing.T) {
	em := getNewEntityManager(t)

	// Create some groups
	if err := em.NewGroup("grp1", "", "", -1); err != nil {
		t.Fatal(err)
	}
	if err := em.NewGroup("grp2", "", "", -1); err != nil {
		t.Fatal(err)
	}
	if err := em.NewGroup("grp3", "", "", -1); err != nil {
		t.Fatal(err)
	}

	// Create some entities
	if err := em.NewEntity("foo", -1, ""); err != nil {
		t.Fatal(err)
	}
	if err := em.NewEntity("bar", -1, ""); err != nil {
		t.Fatal(err)
	}
	if err := em.NewEntity("baz", -1, ""); err != nil {
		t.Fatal(err)
	}

	// Include grp2 and exclude grp3
	if err := em.ModifyGroupExpansions("grp1", "grp2", pb.ExpansionMode_INCLUDE); err != nil {
		t.Fatal(err)
	}
	if err := em.ModifyGroupExpansions("grp1", "grp3", pb.ExpansionMode_EXCLUDE); err != nil {
		t.Fatal(err)
	}

	// Add foo directly, bar via expansion in grp2, and exclude baz via grp3
	if err := em.AddEntityToGroup("foo", "grp1"); err != nil {
		t.Fatal(err)
	}
	if err := em.AddEntityToGroup("bar", "grp2"); err != nil {
		t.Fatal(err)
	}
	if err := em.AddEntityToGroup("baz", "grp1"); err != nil {
		t.Fatal(err)
	}
	if err := em.AddEntityToGroup("baz", "grp3"); err != nil {
		t.Fatal(err)
	}

	members, err := em.ListMembers("grp1")
	if err != nil {
		t.Fatal(err)
	}

	if len(members) != 2 {
		t.Fatal("Wrong membership count")
	}

	if members[0].GetID() != "foo" && members[0].GetID() != "bar" {
		t.Fatal("Bad expansions on grp1", members[0].GetID())
	}
	if members[1].GetID() != "foo" && members[1].GetID() != "bar" {
		t.Fatal("Bad expansions on grp1")
	}
}

func TestAddExpansionInclude(t *testing.T) {
	em := getNewEntityManager(t)

	if err := em.NewGroup("grp1", "", "", -1); err != nil {
		t.Fatal(err)
	}
	if err := em.NewGroup("grp2", "", "", -1); err != nil {
		t.Fatal(err)
	}

	if err := em.ModifyGroupExpansions("grp1", "grp2", pb.ExpansionMode_INCLUDE); err != nil {
		t.Fatal(err)
	}

	g, err := em.GetGroupByName("grp1")
	if err != nil {
		t.Fatal(err)
	}

	if len(g.GetExpansions()) != 1 || g.GetExpansions()[0] != "INCLUDE:grp2" {
		t.Fatal("Expansions incorrect")
	}
}

func TestAddExpansionExclude(t *testing.T) {
	em := getNewEntityManager(t)

	if err := em.NewGroup("grp1", "", "", -1); err != nil {
		t.Fatal(err)
	}
	if err := em.NewGroup("grp2", "", "", -1); err != nil {
		t.Fatal(err)
	}

	if err := em.ModifyGroupExpansions("grp1", "grp2", pb.ExpansionMode_EXCLUDE); err != nil {
		t.Fatal(err)
	}

	g, err := em.GetGroupByName("grp1")
	if err != nil {
		t.Fatal(err)
	}

	if len(g.GetExpansions()) != 1 || g.GetExpansions()[0] != "EXCLUDE:grp2" {
		t.Fatal("Expansions incorrect")
	}
}

func TestDropExpansion(t *testing.T) {
	em := getNewEntityManager(t)

	// Set up some groups
	for _, g := range []string{"grp1", "grp2", "grp3"} {
		if err := em.NewGroup(g, "", "", -1); err != nil {
			t.Fatal(err)
		}
	}

	// Set up some expansions
	if err := em.ModifyGroupExpansions("grp1", "grp2", pb.ExpansionMode_INCLUDE); err != nil {
		t.Fatal(err)
	}
	if err := em.ModifyGroupExpansions("grp1", "grp3", pb.ExpansionMode_INCLUDE); err != nil {
		t.Fatal(err)
	}
	if err := em.ModifyGroupExpansions("grp1", "grp2", pb.ExpansionMode_DROP); err != nil {
		t.Fatal(err)
	}

	// Get the top group
	g, err := em.GetGroupByName("grp1")
	if err != nil {
		t.Fatal(err)
	}

	if len(g.GetExpansions()) != 1 || g.GetExpansions()[0] != "INCLUDE:grp3" {
		t.Fatal("Expansions incorrect")
	}
}

func TestModifyExpansionBadParent(t *testing.T) {
	em := getNewEntityManager(t)

	if err := em.ModifyGroupExpansions("Bogus-Parent", "", pb.ExpansionMode_INCLUDE); err != db.ErrUnknownGroup {
		t.Fatal(err)
	}
}

func TestModifyExpansionBadChild(t *testing.T) {
	em := getNewEntityManager(t)

	if err := em.NewGroup("grp1", "", "", -1); err != nil {
		t.Fatal(err)
	}

	if err := em.ModifyGroupExpansions("grp1", "bogus-child", pb.ExpansionMode_INCLUDE); err != db.ErrUnknownGroup {
		t.Fatal(err)
	}
}

func TestModifyExpansionDuplicate(t *testing.T) {
	em := getNewEntityManager(t)

	// Set up some groups
	for _, g := range []string{"grp1", "grp2"} {
		if err := em.NewGroup(g, "", "", -1); err != nil {
			t.Fatal(err)
		}
	}

	// This should work
	if err := em.ModifyGroupExpansions("grp1", "grp2", pb.ExpansionMode_INCLUDE); err != nil {
		t.Fatal(err)
	}

	// This shouldn't
	if err := em.ModifyGroupExpansions("grp1", "grp2", pb.ExpansionMode_INCLUDE); err != ErrExistingExpansion {
		t.Fatal(err)
	}
}

func TestModifyExpansionCycle(t *testing.T) {
	em := getNewEntityManager(t)

	// Set up some groups
	for _, g := range []string{"grp1", "grp2", "grp3"} {
		if err := em.NewGroup(g, "", "", -1); err != nil {
			t.Fatal(err)
		}
	}

	// This should work
	if err := em.ModifyGroupExpansions("grp1", "grp2", pb.ExpansionMode_INCLUDE); err != nil {
		t.Fatal(err)
	}

	// This should as well
	if err := em.ModifyGroupExpansions("grp2", "grp3", pb.ExpansionMode_INCLUDE); err != nil {
		t.Fatal(err)
	}

	// This one creates the cycle and should fail
	if err := em.ModifyGroupExpansions("grp3", "grp1", pb.ExpansionMode_INCLUDE); err != ErrExistingExpansion {
		t.Fatal(err)
	}
}

func TestCheckGroupCyclesCorruptDB(t *testing.T) {
	em := getNewEntityManager(t)

	// Set up some groups
	for _, g := range []string{"grp1", "grp2", "grp3"} {
		if err := em.NewGroup(g, "", "", -1); err != nil {
			t.Fatal(err)
		}
	}

	if err := em.ModifyGroupExpansions("grp1", "grp2", pb.ExpansionMode_INCLUDE); err != nil {
		t.Fatal(err)
	}

	if err := em.ModifyGroupExpansions("grp2", "grp3", pb.ExpansionMode_INCLUDE); err != nil {
		t.Fatal(err)
	}

	// Make the directory inconsistent by deleting grp2
	if err := em.DeleteGroup("grp2"); err != nil {
		t.Fatal(err)
	}

	// This should bomb out now because there's a cycle check
	// problem loading the group that was deleted.
	if err := em.ModifyGroupExpansions("grp3", "grp1", pb.ExpansionMode_INCLUDE); err != ErrExistingExpansion {
		t.Fatal(err)
	}
}

func TestDedupEntityList(t *testing.T) {
	em := getNewEntityManager(t)

	s := []struct {
		ID     string
		number int32
		secret string
	}{
		{"aaa", -1, ""},
		{"aab", -1, ""},
		{"aac", -1, ""},
		{"aad", -1, ""},
		{"aae", -1, ""},
		{"aaf", -1, ""},
	}

	for _, c := range s {
		if err := em.NewEntity(c.ID, c.number, c.secret); err != nil {
			t.Fatal(err)
		}
	}

	allEntities, err := em.allEntities()
	if err != nil {
		t.Fatal(err)
	}

	// Make a list with duplicates in it
	list := allEntities
	list = append(list, allEntities...)
	if len(list) == len(allEntities) {
		t.Fatal("Lists unexpectedly equal")
	}

	// Dedup it
	list = dedupEntityList(list)

	// Make sure its got no dups
	if len(list) != len(allEntities) {
		t.Fatal("Lists not equal in length")
	}
}

func TestEntityListDifference(t *testing.T) {
	em := getNewEntityManager(t)

	s := []struct {
		ID     string
		number int32
		secret string
	}{
		{"aaa", -1, ""},
		{"aab", -1, ""},
		{"aac", -1, ""},
		{"aad", -1, ""},
		{"aae", -1, ""},
		{"aaf", -1, ""},
	}

	for _, c := range s {
		if err := em.NewEntity(c.ID, c.number, c.secret); err != nil {
			t.Fatal(err)
		}
	}

	// Get a list of everyone
	allEntities, err := em.allEntities()
	if err != nil {
		t.Fatal(err)
	}

	var allButAAA []*pb.Entity
	for _, e := range allEntities {
		if e.GetID() == "aaa" {
			continue
		}
		allButAAA = append(allButAAA, e)
	}

	if len(allButAAA) == len(allEntities) {
		t.Fatal("Lists are not different!")
	}

	shouldJustBeAAA := entityListDifference(allEntities, allButAAA)
	if len(shouldJustBeAAA) != 1 {
		t.Fatalf("Length of shouldJustBeAAA is wrong: %d", len(shouldJustBeAAA))
	}
	if shouldJustBeAAA[0].GetID() != "aaa" {
		t.Fatal("Difference contains wrong result!")
	}
}
