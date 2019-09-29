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

func TestAuthGetToken(t *testing.T) {
	cases := []struct {
		req       pb.AuthRequest
		wantToken string
		wantErr   error
	}{
		{
			req: pb.AuthRequest{
				Entity: &types.Entity{
					ID: proto.String("entity1"),
				},
				Secret: proto.String("secret"),
			},
			wantToken: "{\"EntityID\":\"entity1\",\"Capabilities\":[]}",
			wantErr:   nil,
		},
		{
			req: pb.AuthRequest{
				Entity: &types.Entity{
					ID: proto.String("entity1"),
				},
				Secret: proto.String("wrong-secret"),
			},
			wantToken: "",
			wantErr:   ErrUnauthenticated,
		},
		{
			req: pb.AuthRequest{
				Entity: &types.Entity{
					ID: proto.String("token-issue-error"),
				},
				Secret: proto.String("secret"),
			},
			wantToken: "",
			wantErr:   ErrInternal,
		},
	}

	s := newServer(t)
	initTree(t, s)
	s.CreateEntity("token-issue-error", -1, "secret")

	for i, c := range cases {
		res, err := s.AuthGetToken(context.Background(), &c.req)
		if res.GetToken() != c.wantToken || err != c.wantErr {
			t.Errorf("%d: Want %s %v; Got %s %v", i, c.wantToken, c.wantErr, res.GetToken(), err)
		}
	}
}
