package rpc

import (
	"context"

	pb "github.com/NetAuth/NetAuth/proto"
)

// EntityInfo returns an entity to the client pretty much directly.
// This allows the client to display the various fields of the entity
// without needing to make multiple requests.
func (s *NetAuthServer) EntityInfo(ctx context.Context, netAuthRequest *pb.NetAuthRequest) (*pb.Entity, error) {
	return s.EM.GetEntity(netAuthRequest.GetEntity().GetID())
}
