package rpc2

import (
	"testing"

	types "github.com/NetAuth/Protocol"
)

func TestGetCapabilitiesForEntity(t *testing.T) {
	s := newServer(t)
	initTree(t, s)

	s.CreateGroup("lockout", "", "", -1)
	s.AddEntityToGroup("admin", "lockout")
	s.SetGroupCapability2("lockout", types.Capability_LOCK_ENTITY.Enum())

	caps := s.getCapabilitiesForEntity("admin")
	if len(caps) != 2 {
		t.Error("Not all caps were found", caps)
	}
}
