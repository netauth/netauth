package rpc2

import (
	"context"

	"github.com/netauth/netauth/internal/health"

	types "github.com/netauth/protocol"
	pb "github.com/netauth/protocol/v2"
)

// SystemCapabilities adjusts the capabilities that are on groups by
// default, or if specified directly on an entity.  These capabilities
// only have meaning within NetAuth.
func (s *Server) SystemCapabilities(ctx context.Context, r *pb.CapabilityRequest) (*pb.Empty, error) {
	if s.readonly {
		s.log.Warn("Mutable request in read-only mode!",
			"method", "SystemCapabilities",
			"client", getClientName(ctx),
			"service", getServiceName(ctx),
		)
		return &pb.Empty{}, ErrReadOnly
	}

	// Validate the token and confirm the holder posses
	// GLOBAL_ROOT.  You might wonder why there isn't a capability
	// to assign other capabilities, but then you start going down
	// the rabbit hole and its much more straightforward to just
	// say that you need to be a global superuser to be able to
	// add more capabilities.
	var err error
	ctx, err = s.checkToken(ctx)
	if err != nil {
		return &pb.Empty{}, err
	}
	if err := s.isAuthorized(ctx, types.Capability_GLOBAL_ROOT); err != nil {
		return &pb.Empty{}, err
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
			"client", getClientName(ctx),
			"service", getServiceName(ctx),
		)
		return &pb.Empty{}, ErrMalformedRequest
	}
	if err != nil {
		s.log.Error("Capability Manipulation Error",
			"capability", r.GetCapability(),
			"direct", r.GetDirect(),
			"target", r.GetTarget(),
			"client", getClientName(ctx),
			"service", getServiceName(ctx),
			"error", err,
		)
		return &pb.Empty{}, ErrInternal
	}

	s.log.Info("Capabilities Modified",
		"capability", r.GetCapability(),
		"direct", r.GetDirect(),
		"target", r.GetTarget(),
		"action", r.GetAction(),
		"client", getClientName(ctx),
		"service", getServiceName(ctx),
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
