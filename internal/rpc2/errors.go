package rpc2

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// ErrRequestorUnqualified is returned if the requesting
	// entity does not possess the correct permissions needed to
	// carry out the requested actions.
	ErrRequestorUnqualified = status.Errorf(codes.PermissionDenied, "You do not have permission to carry out that action")

	// ErrMalformedRequest is sent back during some modal requests
	// where the requests has been improperly assembled and cannot
	// be handled at all.
	ErrMalformedRequest = status.Errorf(codes.InvalidArgument, "The request is malformed, consult the protocol documentation and try again")

	// ErrInternal is returned when some backing API has failed to
	// perform as expected.  This is generally for tasks that
	// *should* succeed, but don't for some not automatically
	// detectable error.
	ErrInternal = status.Errorf(codes.Internal, "An internal error has occured and the request could not be processed")

	// ErrUnauthenticated is returned if authentication
	// information cannot be derived, loaded, or validated for a
	// given request.  This is distinct from when authentication
	// information can be derived, but it is insuffucient to
	// perform the requested action.
	ErrUnauthenticated = status.Errorf(codes.Unauthenticated, "Authentication failed")
)
