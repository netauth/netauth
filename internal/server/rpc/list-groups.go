package rpc

import (
	"context"

	pb "github.com/NetAuth/NetAuth/pkg/proto"
)

// ListGroups returns a list to the client.  This is necessary to
// allow the client to display groups for further operations.
func (s *NetAuthServer) ListGroups(ctx context.Context, request *pb.GroupListRequest) (*pb.GroupList, error) {
	groups, err := s.EM.ListGroups()
	if err != nil {
		return nil, err
	}

	reply := new(pb.GroupList)
	reply.Groups = groups

	return reply, nil
}
