package rpc

import (
	"context"
	"fmt"
	"log"

	"github.com/NetAuth/NetAuth/internal/server/entity_manager"

	pb "github.com/NetAuth/NetAuth/proto"
)


// NewEntity takes in a ModEntityRequest and attempts to create a new
// entity.  This call must be made by an entity that has the
// CREATE_ENTITY capability to succeed.
func (s *NetAuthServer) NewEntity(ctx context.Context, modEntityRequest *pb.ModEntityRequest) (*pb.SimpleResult, error) {
	log.Printf("New Entity '%s' was requested by %s",
		modEntityRequest.GetModEntity().GetID(),
		modEntityRequest.GetEntity().GetID())

	result := new(pb.SimpleResult)

	// The requestID and requestSecret are used to authorize the
	// call.
	requestID := modEntityRequest.GetEntity().GetID()
	requestSecret := modEntityRequest.GetEntity().GetSecret()

	// The newID, newUIDNumber, and newSecret are used to populate
	// the corresponding fields on the new entity.  Of these
	// fields, only the newID is strictly required and the others
	// may safely be left at zero values.
	newID := modEntityRequest.GetModEntity().GetID()
	newUIDNumber := modEntityRequest.GetModEntity().GetUidNumber()
	newSecret := modEntityRequest.GetModEntity().GetSecret()

	// After attempting to create the new entity parse out the
	// error message into the string field of the response proto.
	// This needs to be formatted directly for being displayed to
	// a human.
	err := entity_manager.NewEntity(requestID, requestSecret, newID, newUIDNumber, newSecret)
	success := false
	msg := ""
	if err != nil {
		success = false
		msg = fmt.Sprintf("%s", err)
	} else {
		success = true
		msg = "Entity Created"
	}

	result.Success = &success
	result.Msg = &msg

	return result, err
}
