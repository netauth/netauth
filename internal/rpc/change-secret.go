package rpc

import (
	"context"
	"fmt"
	"log"

	pb "github.com/NetAuth/NetAuth/pkg/proto"
)

// ChangeSecret takes in a ModEntityRequest and attempts to change the
// secret on the given entity.  The entity requesting the change must
// either be the one that wishes to change the secret, or must have
// the CHANGE_SECRET capability.
func (s *NetAuthServer) ChangeSecret(ctx context.Context, modEntityRequest *pb.ModEntityRequest) (*pb.SimpleResult, error) {

	log.Printf("Entity '%s' requested secret change for '%s'",
		modEntityRequest.GetEntity().GetID(),
		modEntityRequest.GetModEntity().GetID())

	// The requestID and requestSecret are used to authorize the
	// call.
	requestID := modEntityRequest.GetEntity().GetID()
	requestSecret := modEntityRequest.GetEntity().GetSecret()

	changeID := modEntityRequest.GetModEntity().GetID()
	changeSecret := modEntityRequest.GetModEntity().GetSecret()

	// Change the secret if possible and then return to the client
	// the status.
	err := s.EM.ChangeSecret(requestID, requestSecret, changeID, changeSecret)
	success := false
	msg := ""
	if err != nil {
		success = false
		msg = fmt.Sprintf("%s", err)
	} else {
		success = true
		msg = "Secret Changed"
	}

	result := new(pb.SimpleResult)
	result.Success = &success
	result.Msg = &msg

	return result, err
}
