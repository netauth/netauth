package rpc2

import (
	"context"

	"github.com/NetAuth/NetAuth/internal/token"

	types "github.com/NetAuth/Protocol"
	pb "github.com/NetAuth/Protocol/v2"
)

// AuthEntity handles the process of actually authenticating an
// entity, but does not issue a token.
func (s *Server) AuthEntity(ctx context.Context, r *pb.AuthRequest) (*pb.Empty, error) {
	e := r.GetEntity()
	info := r.GetInfo()

	if err := s.ValidateSecret(e.GetID(), r.GetSecret()); err != nil {
		s.log.Info("Authentication Failed",
			"entity", e.GetID(),
			"service", info.GetService(),
			"client", info.GetID())
		return &pb.Empty{}, ErrUnauthenticated
	}
	s.log.Info("Authentication Succeeded",
		"entity", e.GetID(),
		"service", info.GetService(),
		"client", info.GetID())
	return &pb.Empty{}, nil
}

// AuthGetToken performs entity authentication and issues a token if
// this authentication is successful.
func (s *Server) AuthGetToken(ctx context.Context, r *pb.AuthRequest) (*pb.AuthResult, error) {
	// Check Authentication using the same flow as above.
	_, err := s.AuthEntity(ctx, r)
	if err != nil {
		return &pb.AuthResult{}, err
	}

	caps := s.getCapabilitiesForEntity(*r.Entity.ID)

	// Generate Token
	tkn, err := s.Generate(
		token.Claims{
			EntityID:     r.GetEntity().GetID(),
			Capabilities: caps,
		},
		token.GetConfig(),
	)
	if err != nil {
		s.log.Warn("Error Issuing Token",
			"entity", r.Entity.ID,
			"capabilities", caps,
			"service", r.GetInfo().GetService(),
			"client", r.GetInfo().GetID(),
			"error", err,
		)
		return &pb.AuthResult{}, ErrInternal
	}

	s.log.Info("Token Issued",
		"entity", r.Entity.ID,
		"capabilities", caps,
		"service", r.GetInfo().GetService(),
		"client", r.GetInfo().GetID(),
	)
	return &pb.AuthResult{Token: &tkn}, nil
}

// AuthValidateToken performs server-side verification of a previously
// issued token.  This allows symmetric token algorithms to be used.
func (s *Server) AuthValidateToken(ctx context.Context, r *pb.AuthRequest) (*pb.Empty, error) {
	if _, err := s.Validate(r.GetToken()); err != nil {
		return &pb.Empty{}, ErrUnauthenticated
	}
	return &pb.Empty{}, nil
}

// AuthChangeSecret handles the process of rotating out a stored
// secret for an entity.  This is only appropriate for use in the case
// where NetAuth is maintaining total knowledge of secrets, if this is
// not the case you may need to alter secrets in an external system.
// There are two possible flows depending on if the entity is trying
// to change its own secret or not.  In the first case, the entity
// must be in posession of the original secret, not just a token.  In
// the latter case, the token must have CHANGE_ENTITY_SECRET to
// succeed.
func (s *Server) AuthChangeSecret(ctx context.Context, r *pb.EntityRequest) (*pb.Empty, error) {
	e := r.GetEntity()
	de := r.GetData()
	client := r.GetInfo()

	// While technically a non-local secret database would allow
	// this to proceed, we instead require that mutating requests
	// always hit a fully writeable server.
	if s.readonly {
		s.log.Warn("Mutable request in read-only mode!",
			"method", "AuthChangeSecret",
			"client", client.GetID(),
			"service", client.GetService(),
		)
		return &pb.Empty{}, ErrReadOnly
	}

	// Changes to self require original password, changes to
	// others require CHANGE_ENTITY_SECRET.
	if e.GetID() == de.GetID() {
		if err := s.ValidateSecret(e.GetID(), e.GetSecret()); err != nil {
			s.log.Info("Permission Denied for AuthChangeSecret",
				"modself", true,
				"entity", e.GetID(),
				"service", client.GetService(),
				"client", client.GetID(),
			)
			return &pb.Empty{}, ErrUnauthenticated
		}
	} else {
		c, err := s.Validate(r.GetAuth().GetToken())
		if err != nil {
			s.log.Info("Permission Denied for AuthChangeSecret",
				"modself", false,
				"entity", e.GetID(),
				"service", client.GetService(),
				"client", client.GetID(),
				"error", err,
			)
			return &pb.Empty{}, ErrUnauthenticated
		}
		if !c.HasCapability(types.Capability_CHANGE_ENTITY_SECRET) {
			s.log.Info("Permission Denied for AuthChangeSecret",
				"modself", true,
				"entity", e.GetID(),
				"service", client.GetService(),
				"client", client.GetID(),
				"error", "missing-capability",
			)
			return &pb.Empty{}, ErrRequestorUnqualified
		}
	}

	// Set the secret
	if err := s.SetSecret(de.GetID(), de.GetSecret()); err != nil {
		s.log.Warn("Secret Manipulation Error",
			"entity", de.GetID(),
			"service", client.GetService(),
			"client", client.GetID(),
			"error", err,
		)
		return &pb.Empty{}, ErrInternal
	}
	s.log.Info("Secret Changed",
		"entity", de.GetID(),
		"service", client.GetService(),
		"client", client.GetID(),
	)
	return &pb.Empty{}, nil
}
