package rpc

import (
	"github.com/NetAuth/NetAuth/internal/crypto"
	"github.com/NetAuth/NetAuth/internal/db"
	"github.com/NetAuth/NetAuth/internal/token"
	"github.com/NetAuth/NetAuth/internal/tree"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *NetAuthServer) manageByMembership(entityID, groupName string) bool {
	g, err := s.Tree.GetGroupByName(groupName)
	if err != nil {
		// If the group can't be summoned, pessimistically
		// return false
		return false
	}

	// management by group membership is only available if the
	// group is configured to trust another group for this task,
	// so if this is cleared then no group is trusted.
	if g.GetManagedBy() == "" {
		// This group doesn't have delegated administrative
		// properties.
		return false
	}

	// Get the entity itself for a group check
	e, err := s.Tree.GetEntity(entityID)
	if err != nil {
		return false
	}

	// Always include indirects when evaluating if in an
	// administrative group
	groups := s.Tree.GetMemberships(e, true)

	// Check if any of the groups are the one that grants this
	// power
	for _, name := range groups {
		if name == groupName {
			return true
		}
	}

	// Group checks fall through, return false
	return false
}

// toWireError maps from all of NetAuth's internal errors to canonical
// error codes in gRPC.  This makes interfacing with NetAuth much
// easier for other developers since there is a clear understanding of
// what errors are fatal and what can be retried.  Really the correct
// way to do this is to just have the errors that are returned out
// implement the right type and only intercept the ones that matter in
// here.  That involves doing lots of type checking along the way and
// adding more involved types.  It also risks exposing secure
// implementation back to the client from things kicking up out of
// lower levels.  We'll just use the prepared strings instead.
func toWireError(err error) error {
	switch err {
	case nil:
		return status.Errorf(codes.OK, "Completed sucessfully")
	case crypto.InternalError:
		return status.Errorf(codes.Internal, err.Error())
	case crypto.AuthorizationFailure:
		return status.Errorf(codes.Unauthenticated, err.Error())
	case db.UnknownEntity:
		return status.Errorf(codes.NotFound, err.Error())
	case db.UnknownGroup:
		return status.Errorf(codes.NotFound, err.Error())
	case token.KeyUnavailable:
		return status.Errorf(codes.FailedPrecondition, err.Error())
	case token.TokenInvalid:
		return status.Errorf(codes.Unauthenticated, err.Error())
	case tree.DuplicateEntityID:
		return status.Errorf(codes.AlreadyExists, err.Error())
	case tree.DuplicateGroupName:
		return status.Errorf(codes.AlreadyExists, err.Error())
	case tree.DuplicateNumber:
		return status.Errorf(codes.AlreadyExists, err.Error())
	case tree.UnknownCapability:
		return status.Errorf(codes.NotFound, err.Error())
	case tree.ExistingExpansion:
		return status.Errorf(codes.AlreadyExists, err.Error())
	case MalformedRequest:
		return status.Errorf(codes.InvalidArgument, err.Error())
	case RequestorUnqualified:
		return status.Errorf(codes.PermissionDenied, err.Error())
	case InternalError:
		return status.Errorf(codes.Internal, err.Error())
	default:
		return status.Errorf(codes.Unknown, "An unidentifiable error has occured")
	}
}
