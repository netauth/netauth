package rpc

import (
	"context"
	"fmt"
	"log"

	pb "github.com/NetAuth/NetAuth/pkg/proto"
)

// NewGroup takes in a ModGroupRequest and attempts to create a new
// group.
func (s *NetAuthServer) NewGroup(ctx context.Context, modGroupRequest *pb.ModGroupRequest) (*pb.SimpleResult, error) {
	log.Printf("New Group '%s' was requested by %s",
		modGroupRequest.GetGroup().GetName(),
		modGroupRequest.GetEntity().GetID())

	result := new(pb.SimpleResult)

	// The requestID and requestSecret are used to authorize the
	// call.
	requestID := modGroupRequest.GetEntity().GetID()
	requestSecret := modGroupRequest.GetEntity().GetSecret()

	// The newName, newDisplayName, and newGidNumber are used to
	// populate the new group.
	newName := modGroupRequest.GetGroup().GetName()
	newDisplayName := modGroupRequest.GetGroup().GetDisplayName()
	newGidNumber := modGroupRequest.GetGroup().GetGidNumber()

	err := s.EM.NewGroup(requestID, requestSecret, newName, newDisplayName, newGidNumber)
	success := false
	msg := ""
	if err != nil {
		success = false
		msg = fmt.Sprintf("%s", err)
	} else {
		success = true
		msg = "Group Created"
	}

	result.Success = &success
	result.Msg = &msg

	return result, err
}
