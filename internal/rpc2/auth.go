package rpc2

import (
	"context"

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
	return &pb.AuthResult{}, nil
}

// AuthValidateToken performs server-side verification of a previously
// issued token.  This allows symmetric token algorithms to be used.
func (s *Server) AuthValidateToken(ctx context.Context, r *pb.AuthRequest) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

// AuthChangeSecret handles the process of rotating out a stored
// secret for an entity.  This is only appropriate for use in the case
// where NetAuth is maintaining total knowledge of secrets, if this is
// not the case you may need to alter secrets in an external system.
func (s *Server) AuthChangeSecret(ctx context.Context, r *pb.EntityRequest) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}
