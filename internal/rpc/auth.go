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

	// Get the full fledged entity
	e, err := s.Tree.GetEntity(e.GetID())
	if err != nil {
		log.Println("Entity Vanished!")
	}

	// Get the capabilities list for the token
	capabilities := []string{}
	if e.GetMeta() != nil {
		for _, c := range e.GetMeta().GetCapabilities() {
			capabilities = append(capabilities, pb.Capability_name[int32(c)])
		}
	}

	// Successfully authenticated, now to construct a token
	log.Println(capabilities)
	claims := token.Claims{
		EntityID:     e.GetID(),
		Capabilities: capabilities,
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

func (s *NetAuthServer) ChangeSecret(ctx context.Context, r *pb.ModEntityRequest) (*pb.SimpleResult, error) {
	client := r.GetInfo()
	e := r.GetEntity()
	me := r.GetModEntity()
	t := r.GetAuthToken()

	modself := false

	// Determine if this is a self modifying request or not
	if me != nil && e.GetID() == me.GetID() {
		modself = true
	}

	// Self modifying requests require the original password to
	// proceed
	err := s.Tree.ValidateSecret(e.GetID(), e.GetSecret())
	if modself && err == nil {
		err := s.Tree.SetEntitySecretByID(me.GetID(), me.GetSecret())
		if err != nil {
			return &pb.SimpleResult{
				Success: proto.Bool(false),
				Msg:     proto.String("Error Authenticating"),
			}, nil
		}
		log.Printf("Secret for %s changed (%s@%s)",
			me.GetID(),
			client.GetService(),
			client.GetID())
		return &pb.SimpleResult{
			Success: proto.Bool(true),
			Msg:     proto.String("Secret Changed"),
		}, nil
	}
	if err != nil {
		// Problem with the password in the self modifying
		// request, bail out!
		return &pb.SimpleResult{
			Success: proto.Bool(false),
			Msg:     proto.String("Error Authenticating"),
		}, nil
	}

	// This change is being done administratively since modself
	// was false, so this needs a valid token to proceed.
	c, err := s.Token.Validate(t)
	if err != nil || !c.HasCapability("CHANGE_ENTITY_SECRET") {
		return &pb.SimpleResult{
			Success: proto.Bool(false),
			Msg:     proto.String("Error Authenticating"),
		}, nil
	}

	// Change the secret per what was specified in the
	// modification entity struct.
	if err = s.Tree.SetEntitySecretByID(me.GetID(), me.GetSecret()); err != nil {
		return &pb.SimpleResult{}, err
	}

	// Log this as an administrative change.
	log.Printf("Secret for %s administratively changed by %s (%s@%s)",
		me.GetID(),
		c.EntityID,
		client.GetService(),
		client.GetID())
	return &pb.SimpleResult{
		Success: proto.Bool(true),
		Msg:     proto.String("Secret Changed"),
	}, nil
}
