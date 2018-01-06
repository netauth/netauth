package rpc

import (
	"context"
	"fmt"
	"log"

	"github.com/NetAuth/NetAuth/internal/server/entity_manager"

	pb "github.com/NetAuth/NetAuth/proto"
)

func (s *NetAuthServer) AuthEntity(ctx context.Context, netAuthRequest *pb.NetAuthRequest) (*pb.AuthResult, error) {
	// This must always be defaulted to false here.  Arguably the
	// security of the entire system stems from here where this
	// starts out as false and will require a positive action
	// below to set it true.
	var success = false

	// Go ahead and say who is making this request, and from
	// where, and for what.  This is for diagnostics, and is not
	// really intended to be used for security purposes, but can
	// be nice to look at if things fail below.
	log.Printf("Authenticating %s for service %s to client %s",
		netAuthRequest.GetEntity().GetID(),
		netAuthRequest.GetServiceID(),
		netAuthRequest.GetClientID())

	// Construct and return the response.
	result := new(pb.AuthResult)
	entityID := netAuthRequest.GetEntity().GetID()
	entitySecret := netAuthRequest.GetEntity().GetSecret()
	authStatus := entity_manager.ValidateEntitySecretByID(entityID, entitySecret)
	if authStatus != nil {
		success = false
	} else {
		success = true
	}
	result.Success = &success
	msg := fmt.Sprint(authStatus)
	result.Msg = &msg
	return result, nil
}
