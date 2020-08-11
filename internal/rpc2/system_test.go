package rpc2

import (
	"context"
	"testing"

	"github.com/golang/protobuf/proto"

	types "github.com/netauth/protocol"
	pb "github.com/netauth/protocol/v2"
)

func TestSystemCapabilitiesReadOnly(t *testing.T) {
	s := newServer(t)
	s.readonly = true

	_, err := s.SystemCapabilities(context.Background(), &pb.CapabilityRequest{})
	if err == nil {
		t.Error("Server willing to perform in read-only mode")
	}
}

func TestSsytemCapabilitiesNoAuthentication(t *testing.T) {
	s := newServer(t)

	req := pb.CapabilityRequest{
		Target:     proto.String("group1"),
		Action:     pb.Action_ADD.Enum(),
		Capability: types.Capability_CREATE_ENTITY.Enum(),
	}

	_, err := s.SystemCapabilities(InvalidAuthContext, &req)
	if err == nil {
		t.Log(err)
		t.Error("Request with invalidated token was accepted")
	}
}

func TestSsytemCapabilitiesNoCapability(t *testing.T) {
	s := newServer(t)

	req := pb.CapabilityRequest{
		Target:     proto.String("group1"),
		Action:     pb.Action_ADD.Enum(),
		Capability: types.Capability_CREATE_ENTITY.Enum(),
	}

	_, err := s.SystemCapabilities(UnprivilegedContext, &req)
	if err == nil {
		t.Log(err)
		t.Error("Request with invalidated token was accepted")
	}
}

func TestSystemCapabilitiesEntity(t *testing.T) {
	s := newServer(t)
	initTree(t, s.Manager)

	req := pb.CapabilityRequest{
		Target:     proto.String("entity1"),
		Direct:     proto.Bool(true),
		Action:     pb.Action_ADD.Enum(),
		Capability: types.Capability_CREATE_ENTITY.Enum(),
	}

	_, err := s.SystemCapabilities(PrivilegedContext, &req)
	if err != nil {
		t.Log(err)
		t.Fatal("Request with validated token was rejected")
	}

	e, _ := s.FetchEntity("entity1")
	if e.GetMeta() == nil {
		t.Fatal("Failed to set capability on entity (no meta)")
	}
	if len(e.Meta.Capabilities) != 1 || e.Meta.Capabilities[0] != types.Capability_CREATE_ENTITY {
		t.Error("Failed to set capability on entity")
	}

	req.Action = pb.Action_DROP.Enum()
	_, err = s.SystemCapabilities(PrivilegedContext, &req)
	if err != nil {
		t.Log(err)
		t.Error("Request with validated token was rejected")
	}

	e, _ = s.FetchEntity("entity1")
	if len(e.Meta.Capabilities) != 0 {
		t.Error("Failed to remove capability from entity")
	}

}

func TestSystemCapabilitiesGroup(t *testing.T) {
	s := newServer(t)
	initTree(t, s.Manager)

	req := pb.CapabilityRequest{
		Target:     proto.String("group1"),
		Action:     pb.Action_ADD.Enum(),
		Capability: types.Capability_GLOBAL_ROOT.Enum(),
	}

	_, err := s.SystemCapabilities(PrivilegedContext, &req)
	if err != nil {
		t.Log(err)
		t.Error("Request with validated token was rejected")
	}

	g, _ := s.FetchGroup("group1")
	if len(g.Capabilities) != 1 || g.Capabilities[0] != types.Capability_GLOBAL_ROOT {
		t.Error("Failed to set capability on group")
	}

	req.Action = pb.Action_DROP.Enum()
	_, err = s.SystemCapabilities(PrivilegedContext, &req)
	if err != nil {
		t.Log(err)
		t.Error("Request with validated token was rejected")
	}

	g, _ = s.FetchGroup("group1")
	if len(g.Capabilities) != 0 {
		t.Error("Failed to remove capability from group")
	}

}

func TestSystemCapabilitiesMalformedRequest(t *testing.T) {
	s := newServer(t)

	_, err := s.SystemCapabilities(PrivilegedContext, &pb.CapabilityRequest{})
	if err != ErrMalformedRequest {
		t.Log(err)
		t.Error("Request with invalidated token was accepted")
	}
}

func TestSystemCapabilitiesManipulationError(t *testing.T) {
	s := newServer(t)
	initTree(t, s.Manager)

	req := pb.CapabilityRequest{
		Direct: proto.Bool(true),
		Target: proto.String("entity1"),
	}

	_, err := s.SystemCapabilities(PrivilegedContext, &req)
	if err != ErrInternal {
		t.Log(req.Capability)
		t.Log(err)
		t.Error("Malformed request returned no error")
	}
}

func TestSystemPing(t *testing.T) {
	s := newServer(t)

	_, err := s.SystemPing(context.Background(), &pb.Empty{})
	if err != nil {
		t.Error("Non nil error", err)
	}
}

func TestSystemStatus(t *testing.T) {
	s := newServer(t)

	status, err := s.SystemStatus(context.Background(), &pb.Empty{})
	if err != nil {
		t.Error("Non nil error", err)
	}

	if status.GetSystemOK() != true {
		t.Error("Status does not reflect green state")
	}
}
