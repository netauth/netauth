package rpc2

import (
	"context"

	"github.com/NetAuth/NetAuth/internal/health"

	pb "github.com/NetAuth/Protocol/v2"
)

// SystemCapabilities adjusts the capabilities that are on groups by
// default, or if specified directly on an entity.  These capabilities
// only have meaning within NetAuth.
func (s *Server) SystemCapabilities(ctx context.Context, r *pb.CapabilityRequest) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

// SystemPing provides the most simple "the server is alive" check.
// It does not provide any additional information, if you want that
// use SystemStatus.
func (s *Server) SystemPing(ctx context.Context, r *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

// SystemStatus returns detailed status information on the server.
func (s *Server) SystemStatus(ctx context.Context, r *pb.Empty) (*pb.ServerStatus, error) {
	status := health.Check()
	return status.Proto(), nil
}
