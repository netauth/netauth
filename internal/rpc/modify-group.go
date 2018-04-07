package rpc

import (
	"context"
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"

	pb "github.com/NetAuth/NetAuth/pkg/proto"
)

// ModifyGroupMeta gets the group, extracts the credentials and update
// group value and attempts to modify the group in place.
func (s *NetAuthServer) ModifyGroupMeta(ctx context.Context, modGroupRequest *pb.ModGroupRequest) (*pb.SimpleResult, error) {

	log.Printf("Entity '%s' requested metadata update for '%s'",
		modGroupRequest.GetEntity().GetID(),
		modGroupRequest.GetGroup().GetName())

	// The requestID and requestSecret are used to authorize the
	// call.
	requestID := modGroupRequest.GetEntity().GetID()
	requestSecret := modGroupRequest.GetEntity().GetSecret()

	g := modGroupRequest.GetGroup()

	err := s.EM.UpdateGroupMeta(requestID, requestSecret, g.GetName(), g)

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
