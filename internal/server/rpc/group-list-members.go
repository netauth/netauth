package rpc

import (
	"context"
	"log"

	pb "github.com/NetAuth/NetAuth/proto"
)

// ListGroupMembers takes in a Group and lists the members that are
// contained by that group.  No distinction is made between direct and
// indirect members.
func (s *NetAuthServer) ListGroupMembers(ctx context.Context, groupRequest *pb.GroupMemberRequest) (*pb.EntityList, error) {
	log.Printf("Group members list request for '%s' (service: %s; client %s)",
		*groupRequest.Group.Name,
		*groupRequest.ServiceID,
		*groupRequest.ClientID)

	entityList := new(pb.EntityList)
	members, err := s.EM.ListMembers(*groupRequest.Group.Name)
	entityList.Members = members
	return entityList, err
}
