package mresolver

import (
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

func testAtom(mr *MResolver) {
	mr.SyncDirectGroups("entity1", []string{"group1"})
	mr.SyncDirectGroups("entity2", []string{"group2", "group4"})

	// This configures the group caches, its very important that
	// this is manually checked for cycles when editing, as there
	// are no other checks at this layer.
	mr.atom.gc = map[string]*resolvableGroup{
		"group0": {
			// group0 is only ever resolved implicitly,
			// which checks to make sure that fallthrough
			// resolution works correctly.
			self:    "group0",
			exclude: []string{"group2"},
		},
		"group1": {
			self: "group1",
		},
		"group2": {
			self:    "group2",
			include: []string{"group1"},
		},
		"group3": {
			self:    "group3",
			exclude: []string{"group2"},
		},
		"group4": {
			self:    "group4",
			include: []string{"group1"},
			exclude: []string{"group3"},
		},
		"group5": {
			self:    "group5",
			include: []string{"group2"},
			exclude: []string{"group0"},
		},
	}

	// This is only partial to test that removing works.
	mr.atom.ga = map[string]map[string]struct{}{
		"group0": {
			"group5": {},
		},
		"group1": {
			"group2": {},
		},
	}
}

func TestSyncDirectGroups(t *testing.T) {
	x := New()
	assert.Equal(t, 0, len(x.atom.dm))
	x.SyncDirectGroups("entity1", []string{"group1", "group2"})
	assert.Equal(t, 1, len(x.atom.dm))
	assert.Equal(t, 2, len(x.atom.dm["entity1"]))

	// Now delete it
	x.RemoveEntity("entity1")
	assert.Equal(t, 0, len(x.atom.dm))
}

func TestSyncGroup(t *testing.T) {
	x := New()
	l := hclog.New(&hclog.LoggerOptions{Name: "debug"})
	l.SetLevel(hclog.Trace)
	x.SetParentLogger(l)

	// Setup a chain of groups
	x.SyncGroup("group1", []string{}, []string{})
	x.SyncGroup("group2", []string{"group1"}, []string{})
	x.SyncGroup("group3", []string{"group2"}, []string{})
	x.SyncGroup("group4", []string{"group3"}, []string{})
	x.SyncGroup("group5", []string{"group4"}, []string{})
	x.SyncGroup("group6", []string{}, []string{"group5"})
	assert.Equal(t, "(group5|(group4|(group3|(group2|group1))))", x.atom.gr["group5"].String())
	assert.Equal(t, "(group6&!(group5|(group4|(group3|(group2|group1)))))", x.atom.gr["group6"].String())

	// Change group 2 to not include group1 anymore and make sure
	// that group5 updates.
	x.SyncGroup("group2", []string{}, []string{})
	assert.Equal(t, "(group5|(group4|(group3|group2)))", x.atom.gr["group5"].String())
	assert.Equal(t, "(group6&!(group5|(group4|(group3|group2))))", x.atom.gr["group6"].String())
}

func TestRemoveGroup(t *testing.T) {
	x := New()
	testAtom(x)
	x.RemoveGroup("group3")
	_, ok := x.atom.gc["group3"]
	assert.Equal(t, false, ok)
}

func TestResolve(t *testing.T) {
	x := New()
	testAtom(x)

	// This order is very important, it ensures that there is a
	// cache miss on include and another one on exclude.
	x.Resolve("group1")
	x.Resolve("group5")
	x.Resolve("group4")

	check := map[string]string{
		"group1": "group1",
		"group4": "((group4|group1)&!(group3&!(group2|group1)))",
		"group5": "((group5|(group2|group1))&!(group0&!(group2|group1)))",
	}

	for group, expression := range check {
		assert.Equalf(t, expression, x.atom.gr[group].String(), "Group %s wrong expression", group)
	}
}

func TestResolveBadGroups(t *testing.T) {
	x := New()
	x.atom.gc = map[string]*resolvableGroup{
		"bad-include": {
			self:    "bad-include",
			include: []string{"does-not-exist"},
		},

		"bad-exclude": {
			self:    "bad-exclude",
			exclude: []string{"does-not-exist"},
		},
	}

	assert.Equal(t, ErrInsufficientKnowledge, x.Resolve("does-not-exist"))
	assert.Equal(t, ErrInsufficientKnowledge, x.Resolve("bad-include"))
	assert.Equal(t, ErrInsufficientKnowledge, x.Resolve("bad-exclude"))
}

func TestMembersOfGroup(t *testing.T) {
	x := New()
	testAtom(x)
	for _, g := range []string{"group0", "group1", "group2", "group3", "group4", "group5"} {
		x.Resolve(g)
	}

	assert.EqualValues(t, []string{"entity1"}, x.MembersOfGroup("group1"))
	assert.Equal(t, []string{}, x.MembersOfGroup("does-not-exist"))
}

func TestGroupsForEntity(t *testing.T) {
	x := New()
	testAtom(x)
	for _, g := range []string{"group0", "group1", "group2", "group3", "group4", "group5"} {
		x.Resolve(g)
	}

	assert.ElementsMatch(t, []string{"group2", "group4", "group5", "group1"}, x.GroupsForEntity("entity1"))
	assert.Equal(t, []string{}, x.GroupsForEntity("does-not-exist"))
}
