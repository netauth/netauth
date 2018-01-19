package rpc

import (
	"context"
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
	
	pb "github.com/NetAuth/NetAuth/proto"
)

// ModifyEntityMeta takes in a ModEntityRequest and extracts the
// EntityMeta from the ModEntity.  The requesting entity must either
// be the one being modified or posses the MODIFY_ENTITY_META
// capability.
func (s *NetAuthServer) ModifyEntityMeta(ctx context.Context, modEntityRequest *pb.ModEntityRequest) (*pb.SimpleResult, error) {

	log.Printf("Entity '%s' requested metadata update for '%s'",
		modEntityRequest.GetEntity().GetID(),
		modEntityRequest.GetModEntity().GetID())

	// The requestID and requestSecret are used to authorize the
	// call.
	requestID := modEntityRequest.GetEntity().GetID()
	requestSecret := modEntityRequest.GetEntity().GetSecret()

	modEntityID := modEntityRequest.GetModEntity().GetID()
	modMeta := modEntityRequest.GetModEntity().GetMeta()

	err := s.EM.UpdateEntityMeta(requestID, requestSecret, modEntityID, modMeta)
	
	result := new(pb.SimpleResult)

	if err != nil {
		result.Msg = proto.String(fmt.Sprintf("%s", err))
		result.Success = proto.Bool(false)
	} else {
		result.Msg = proto.String("Successfully updated metadata")
		result.Success = proto.Bool(true)
	}

	return result, nil
}
