package rpc2

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/netauth/netauth/internal/token"

	types "github.com/netauth/protocol"
)

type claimsContextKey struct{}

func (s *Server) getCapabilitiesForEntity(ctx context.Context, id string) []types.Capability {
	// Get the full fledged entity; we can assert no error here
	// since the entity was just loaded to perform an
	// authentication check prior to calling this function.
	// There's a minimal risk that the entity has vanished, but if
	// that's the case then this will return empty since the
	// FetchEntity function will return an empty entity, thus no
	// authorization attack can be performed via this vector.
	e, _ := s.FetchEntity(ctx, id)

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
	groupNames := s.GetMemberships(ctx, e)
	for _, name := range groupNames {
		g, _ := s.FetchGroup(ctx, name)
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
// and used here.  The end result is that claims are added to a
// returned context.  Actually using these claims should be done by
// isAuthorized.
func (s *Server) checkToken(ctx context.Context) (context.Context, error) {
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
		return ctx, ErrMalformedRequest
	}
	c, err := s.Validate(tkn)
	if err != nil {
		s.log.Info("Permission Denied",
			"method", method,
			"client", getClientName(ctx),
			"service", getServiceName(ctx),
			"error", err,
		)
		return ctx, ErrUnauthenticated
	}
	ctx = context.WithValue(ctx, claimsContextKey{}, c)
	return ctx, nil
}

// isAuthorized checks for a specific capability in the claims from
// the context.  If it is not present, then the client is not
// sufficientlly empowered by capabilities alone to make the given
// request, but may be authorized by group membership.
func (s *Server) isAuthorized(ctx context.Context, reqCap types.Capability) error {
	method, ok := grpc.Method(ctx)
	if !ok {
		method = "UNKNOWN"
	}
	c := getTokenClaims(ctx)
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

// getTokenClaims returns the claims from the context without
// modifying it.  The claims will either be populated if a token was
// previously parsed into the context, or empty if no such token has
// been successfully parsed.
func getTokenClaims(ctx context.Context) token.Claims {
	v, ok := ctx.Value(claimsContextKey{}).(token.Claims)
	if !ok {
		return token.Claims{}
	}
	return v
}

// manageByMembership checks if the entity identified by entityID is a
// member of any group that group g has delegated management authority
// to.  In this way, an entity can be allowed to alter certain groups
// without needing to grant broad server level authority.
func (s *Server) manageByMembership(ctx context.Context, entityID string, g *types.Group) bool {
	g, err := s.FetchGroup(ctx, g.GetName())
	if err != nil {
		return false
	}

	// Management by membership is only available if explicitly
	// enabled.  If the value of the string is empty, this group
	// has not delegated management authority to any other groups.
	if g.GetManagedBy() == "" {
		return false
	}

	e, err := s.FetchEntity(ctx, entityID)
	if err != nil {
		return false
	}

	for _, name := range s.GetMemberships(ctx, e) {
		if name == g.GetManagedBy() {
			return true
		}
	}
	return false
}

// mutablePrequisitesAreMet checks for common mutable prerequisites
// such as the server being in a writeable mode, and the correct
// capability being present in a valid token.
func (s *Server) mutablePrequisitesMet(ctx context.Context, c types.Capability) error {
	if s.readonly {
		s.log.Warn("Mutable request in read-only mode!",
			"method", "EntityUM",
			"client", getClientName(ctx),
			"service", getServiceName(ctx),
		)
		return ErrReadOnly
	}

	// Token validation and authorization
	var err error
	ctx, err = s.checkToken(ctx)
	if err != nil {
		return err
	}
	if err := s.isAuthorized(ctx, c); err != nil {
		return err
	}
	return nil
}
