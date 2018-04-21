package rpc

import (
	"context"
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/internal/token"

	pb "github.com/NetAuth/NetAuth/pkg/proto"
)

func (s *NetAuthServer) AuthEntity(ctx context.Context, r *pb.NetAuthRequest) (*pb.SimpleResult, error) {
	// Log out some useful stuff for debugging
	client := r.GetInfo()
	entity := r.GetEntity()
	log.Printf("Authenticating %s (%s@%s)",
		entity.GetID(),
		client.GetService(),
		client.GetID())

	// Construct and return the response.
	result := new(pb.SimpleResult)

	if err := s.Tree.ValidateSecret(entity.GetID(), entity.GetSecret()); err != nil {
		return nil, err
	}

	result.Success = proto.Bool(true)
	result.Msg = proto.String("Entity authentication succeeded")
	return result, nil
}

func (s *NetAuthServer) GetToken(ctx context.Context, r *pb.NetAuthRequest) (*pb.TokenResult, error) {
	client := r.GetInfo()
	e := r.GetEntity()

	log.Printf("Token requested for %s (%s@%s)",
		e.GetID(),
		client.GetService(),
		client.GetID())

	// Run the normal authentication flow
	if err := s.Tree.ValidateSecret(e.GetID(), e.GetSecret()); err != nil {
		return nil, err
	}

	// Successfully authenticated, now to construct a token
	claims := token.Claims{
		EntityID: e.GetID(),
	}

	// Generate the token with the specified claims
	tkn, err := s.Token.Generate(claims, token.GetConfig())
	if err != nil {
		return nil, err
	}

	// Construct the reply containing the token
	reply := pb.TokenResult{
		Success: proto.Bool(true),
		Msg:     proto.String("Token Granted"),
		Token:   &tkn,
	}

	return &reply, nil
}

func (s *NetAuthServer) ValidateToken(ctx context.Context, r *pb.NetAuthRequest) (*pb.SimpleResult, error) {
	client := r.GetInfo()
	e := r.GetEntity()
	t := r.GetAuthToken()

	log.Printf("Token validation requested by %s (%s@%s)",
		e.GetID(),
		client.GetService(),
		client.GetID())

	// These will get cleared to error values if there's a fault
	msg := "Token validation successful"
	success := true

	// Validate the token and if it validates, return that the
	// token is valid.
	_, err := s.Token.Validate(t)
	if err != nil {
		msg = fmt.Sprintf("%s", err)
		success = false
	}

	// Compose and return the reply
	reply := pb.SimpleResult{
		Success: &success,
		Msg:     &msg,
	}
	return &reply, nil
}
