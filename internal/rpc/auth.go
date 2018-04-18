package rpc

import (
	"context"
	"log"

	pb "github.com/NetAuth/NetAuth/pkg/proto"
)

func (s *NetAuthServer) AuthEntity(ctx context.Context, netAuthRequest *pb.NetAuthRequest) (*pb.SimpleResult, error) {
	// This must always be defaulted to false here.  Arguably the
	// security of the entire system stems from here where this
	// starts out as false and will require a positive action
	// below to set it true.
	var success = false

	// Go ahead and say who is making this request, and from
	// where, and for what.  This is for diagnostics, and is not
	// really intended to be used for security purposes, but can
	// be nice to look at if things fail below.
	client := netAuthRequest.GetInfo()
	log.Printf("Authenticating %s for service %s to client %s",
		netAuthRequest.GetEntity().GetID(),
		client.GetID(),
		client.GetService())

	// Construct and return the response.
	result := new(pb.SimpleResult)
	entityID := netAuthRequest.GetEntity().GetID()
	entitySecret := netAuthRequest.GetEntity().GetSecret()
	authStatus := s.Tree.ValidateSecret(entityID, entitySecret)
	msg := ""
	if authStatus != nil {
		success = false
		msg = "Entity authentication failed"
	} else {
		success = true
		msg = "Entity authentication succeeded"
	}
	result.Success = &success
	result.Msg = &msg
	return result, nil
}
