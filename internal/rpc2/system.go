package rpc2

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/NetAuth/NetAuth/internal/health"

	pb "github.com/NetAuth/Protocol/v2"
)

// SystemCapabilities adjusts the capabilities that are on groups by
// default, or if specified directly on an entity.  These capabilities
// only have meaning within NetAuth.
func (s *Server) SystemCapabilities(ctx context.Context, r *pb.CapabilityRequest) (*pb.Empty, error) {
	authdata := r.GetAuth()
	client := r.GetInfo()

	if s.readonly {
		s.log.Warn("Mutable request in read-only mode!",
			"method", "SystemCapabilities",
			"client", client.GetID(),
			"service", client.GetService(),
		)
		return &pb.Empty{}, status.Errorf(codes.FailedPrecondition, "Server is in read-only mode")
	}

	// Validate the token and confirm the holder posses
	// GLOBAL_ROOT.  You might wonder why there isn't a capability
	// to assign other capabilities, but then you start going down
	// the rabbit hole and its much more straightforward to just
	// say that you need to be a global superuser to be able to
	// add more capabilities.
	c, err := s.Validate(authdata.GetToken())
	if err != nil || !c.HasCapability("GLOBAL_ROOT") {
		return nil, ErrRequestorUnqualified
	}

	switch {
	case r.GetDirect() && r.GetAction() == pb.Action_ADD && r.GetTarget() != "":
		err = s.SetEntityCapability2(r.GetTarget(), r.Capability)
	case r.GetDirect() && r.GetAction() == pb.Action_DROP && r.GetTarget() != "":
		err = s.DropEntityCapability2(r.GetTarget(), r.Capability)
	case !r.GetDirect() && r.GetAction() == pb.Action_ADD && r.GetTarget() != "":
		err = s.SetGroupCapability2(r.GetTarget(), r.Capability)
	case !r.GetDirect() && r.GetAction() == pb.Action_DROP && r.GetTarget() != "":
		err = s.DropGroupCapability2(r.GetTarget(), r.Capability)
	default:
		s.log.Warn("Malformed request",
			"method", "SystemCapabilities",
			"client", client.GetID(),
			"service", client.GetService(),
		)
		return &pb.Empty{}, ErrMalformedRequest
	}
	if err != nil {
		s.log.Error("Capability Manipulation Error",
			"capability", r.GetCapability(),
			"direct", r.GetDirect(),
			"target", r.GetTarget(),
			"client", client.GetID(),
			"service", client.GetService(),
			"error", err,
		)
		return &pb.Empty{}, ErrInternal
	}

	s.log.Info("Capabilities Modified",
		"capability", r.GetCapability(),
		"direct", r.GetDirect(),
		"target", r.GetTarget(),
		"action", r.GetAction(),
		"client", client.GetID(),
		"service", client.GetService(),
	)
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
