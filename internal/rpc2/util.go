package rpc2

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	types "github.com/NetAuth/Protocol"
)

func (s *Server) getCapabilitiesForEntity(id string) []types.Capability {
	// Get the full fledged entity; we can assert no error here
	// since the entity was just loaded to perform an
	// authentication check prior to calling this function.
	// There's a minimal risk that the entity has vanished, but if
	// that's the case then this will return empty since the
	// FetchEntity function will return an empty entity, thus no
	// authorization attack can be performed via this vector.
	e, _ := s.FetchEntity(id)

	// First get the capabilities that are provided by the entity
	// itself.
	caps := make(map[types.Capability]int)
	if e.GetMeta() != nil {
		for _, c := range e.GetMeta().GetCapabilities() {
			caps[c]++
		}
	}

	// Next get the capabilities that are provided by any groups
	// the entity may be in; include indirects for authentication
	// queries.  We can assert no error here because the worst
	// case is a group is returned that is completely empty.  In
	// this case there will simply be no additional capabilities
	// granted, and this function can continue without incident.
	groupNames := s.GetMemberships(e, true)
	for _, name := range groupNames {
		g, _ := s.FetchGroup(name)
		for _, c := range g.GetCapabilities() {
			caps[c]++
		}
	}

	// Flatten the capabilities out into a list
	capabilities := make([]types.Capability, len(caps))
	i := 0
	for c := range caps {
		capabilities[i] = c
		i++
	}
	return capabilities
}

// checkToken is used to validate authorization from the context.
// This authorization is present in the form of a token in the
// "authorization" field of the request metadata which is extracted
// and used here.
func (s *Server) checkToken(ctx context.Context, reqCap types.Capability) error {
	tkn := getSingleStringFromMetadata(ctx, "authorization")
	method, ok := grpc.Method(ctx)
	if !ok {
		method = "UNKNOWN"
	}
	if tkn == "" {
		s.log.Info("Request contains no token but token is required!",
			"method", method,
			"client", getClientName(ctx),
			"service", getServiceName(ctx),
		)
		return ErrMalformedRequest
	}
	c, err := s.Validate(tkn)
	if err != nil {
		s.log.Info("Permission Denied",
			"method", method,
			"client", getClientName(ctx),
			"service", getServiceName(ctx),
			"error", err,
		)
		return ErrUnauthenticated
	}

	if !c.HasCapability(reqCap) {
		s.log.Info("Permission Denied",
			"method", method,
			"client", getClientName(ctx),
			"service", getServiceName(ctx),
			"error", "missing-capability",
		)
		return ErrRequestorUnqualified
	}
	return nil
}

// getSingleStringFromMetadata is a convenience function that helps to
// pull individual values from the request metadata.  It asserts that
// only a single value will be set, and that that value is a string.
func getSingleStringFromMetadata(ctx context.Context, key string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	sl := md.Get(key)
	if len(sl) != 1 {
		return ""
	}
	return sl[0]
}

// getClientName returns the client name.  If no name was set, the
// string "BOGUS_CLIENT" is returned.
func getClientName(ctx context.Context) string {
	s := getSingleStringFromMetadata(ctx, "client-name")
	if s == "" {
		return "BOGUS_CLIENT"
	}
	return s
}

// getServiceName returns the service name.  If no name was set the
// string "BOGUS_SERVICE" is returned.
func getServiceName(ctx context.Context) string {
	s := getSingleStringFromMetadata(ctx, "service-name")
	if s == "" {
		return "BOGUS_SERVICE"
	}
	return s
}
