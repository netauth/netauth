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
	ErrInternal = status.Errorf(codes.Internal, "An internal error has occurred and the request could not be processed")

	// ErrUnauthenticated is returned if authentication
	// information cannot be derived, loaded, or validated for a
	// given request.  This is distinct from when authentication
	// information can be derived, but it is insufficient to
	// perform the requested action.
	ErrUnauthenticated = status.Errorf(codes.Unauthenticated, "Authentication failed")

	// ErrReadOnly is returned if the server is in read-only mode
	// and a mutating request is received.  In this case the
	// server cannot comply, and the behavior cannot be retried,
	// so we return that the feature is unimplemented as in this
	// node it might as well be.
	ErrReadOnly = status.Errorf(codes.Unimplemented, "Server is in read-only mode")

	// ErrExists iis returned when creation would create a
	// duplicate resource and this is not handled internally via
	// automatic deduplication.  Examples include trying to create
	// an entity with an existing ID, or a group with an already
	// used number.
	ErrExists = status.Errorf(codes.AlreadyExists, "One or more parameters collides with an existing item")

	// ErrDoesNotExist is, as the name would imply, returned if an
	// action calls for a resource that does not exist.  This can
	// be the case when an update or change is requested on an
	// entity or group that does not exist, or when an expansion
	// that doesn't exist is modified.
	ErrDoesNotExist = status.Errorf(codes.NotFound, "The requested resource does not exist")
)
