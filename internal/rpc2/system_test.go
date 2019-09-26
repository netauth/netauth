package rpc2

import (
	"context"
	"testing"

	pb "github.com/NetAuth/Protocol/v2"
)

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
	t.Logf("Full status: %v", status)
}
