package rpc

import (
	"context"
	"fmt"
	"log"

	pb "github.com/NetAuth/NetAuth/pkg/proto"
)

// DeleteGroup takes in a ModGroupRequest and attempts to delete an
// existing group.
func (s *NetAuthServer) DeleteGroup(ctx context.Context, modGroupRequest *pb.ModGroupRequest) (*pb.SimpleResult, error) {
	log.Printf("Delete of Group '%s' was requested by %s",
		modGroupRequest.GetGroup().GetName(),
		modGroupRequest.GetEntity().GetID())

	result := new(pb.SimpleResult)

	// The requestID and requestSecret are used to authorize the
	// call.
	requestID := modGroupRequest.GetEntity().GetID()
	requestSecret := modGroupRequest.GetEntity().GetSecret()

	// The name is the only thing needed to delete the group, no
	// other fields set matter here.
	name := modGroupRequest.GetGroup().GetName()

	err := s.EM.DeleteGroup(requestID, requestSecret, name)
	success := false
	msg := ""
	if err != nil {
		success = false
		msg = fmt.Sprintf("%s", err)
	} else {
		success = true
		msg = "Group Deleted"
	}

	result.Success = &success
	result.Msg = &msg

	return result, err
}
