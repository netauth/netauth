package rpc

import (
	"context"
	"fmt"
	"log"

	pb "github.com/NetAuth/NetAuth/proto"
)

// RemoveEntity takes in a ModEntityRequest and attempts to remove an existing
// entity.  This call must be made by an entity that has the
// DESTROY_ENTITY capability to succeed.
func (s *NetAuthServer) RemoveEntity(ctx context.Context, modEntityRequest *pb.ModEntityRequest) (*pb.SimpleResult, error) {
	log.Printf("Delete entity '%s' was requested by '%s'",
		modEntityRequest.GetModEntity().GetID(),
		modEntityRequest.GetEntity().GetID())

	result := new(pb.SimpleResult)

	// The requestID and requestSecret are used to authorize the
	// call.
	requestID := modEntityRequest.GetEntity().GetID()
	requestSecret := modEntityRequest.GetEntity().GetSecret()

	delID := modEntityRequest.GetModEntity().GetID()

	// After attempting to delete the entity, parse out the error
	// and return it to the RPC layer.
	err := s.EM.DeleteEntity(requestID, requestSecret, delID)
	success := false
	msg := ""
	if err != nil {
		success = false
		msg = fmt.Sprintf("%s", err)
	} else {
		success = true
		msg = "Entity Removed"
	}

	result.Success = &success
	result.Msg = &msg

	return result, err
}
