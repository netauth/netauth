package rpc2

import (
	"context"
	"testing"

	"github.com/golang/protobuf/proto"

	types "github.com/NetAuth/Protocol"
	pb "github.com/NetAuth/Protocol/v2"
)

func TestAuthEntity(t *testing.T) {
	cases := []struct {
		req     pb.AuthRequest
		wantErr error
	}{
		{
			req: pb.AuthRequest{
				Entity: &types.Entity{
					ID: proto.String("entity1"),
				},
				Secret: proto.String("secret"),
			},
			wantErr: nil,
		},
		{
			req: pb.AuthRequest{
				Entity: &types.Entity{
					ID: proto.String("entity1"),
				},
				Secret: proto.String("wrong"),
			},
			wantErr: ErrUnauthenticated,
		},
	}

	s := newServer(t)
	initTree(t, s)
	for i, c := range cases {
		if _, err := s.AuthEntity(context.Background(), &c.req); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}
