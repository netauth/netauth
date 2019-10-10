package rpc2

import (
	"context"
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/internal/token/null"

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

func TestAuthValidateToken(t *testing.T) {
	cases := []struct {
		token   string
		wantErr error
	}{
		{null.ValidToken, nil},
		{null.InvalidToken, ErrUnauthenticated},
	}

	req := pb.AuthRequest{}
	s := newServer(t)
	for i, c := range cases {
		req.Token = &c.token
		if _, err := s.AuthValidateToken(context.Background(), &req); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}

func TestAuthChangeSecret(t *testing.T) {
	cases := []struct {
		ctx      context.Context
		req      pb.EntityRequest
		readonly bool
		wantErr  error
	}{
		{
			// Works, and changes own secret
			ctx: UnauthenticatedContext,
			req: pb.EntityRequest{
				Entity: &types.Entity{
					ID:     proto.String("entity1"),
					Secret: proto.String("secret"),
				},
				Data: &types.Entity{
					ID:     proto.String("entity1"),
					Secret: proto.String("secret1"),
				},
			},
			readonly: false,
			wantErr:  nil,
		},
		{
			// Fails, original secret not available
			ctx: UnauthenticatedContext,
			req: pb.EntityRequest{
				Entity: &types.Entity{
					ID:     proto.String("entity1"),
					Secret: proto.String("incorrect-secret"),
				},
				Data: &types.Entity{
					ID:     proto.String("entity1"),
					Secret: proto.String("secret1"),
				},
			},
			readonly: false,
			wantErr:  ErrUnauthenticated,
		},
		{
			// Fails, read-only
			ctx: UnauthenticatedContext,
			req: pb.EntityRequest{
				Entity: &types.Entity{
					ID:     proto.String("entity1"),
					Secret: proto.String("secret"),
				},
				Data: &types.Entity{
					ID:     proto.String("entity1"),
					Secret: proto.String("secret1"),
				},
			},
			readonly: true,
			wantErr:  ErrReadOnly,
		},
		{
			// Works, auth'd via token
			ctx: PrivilegedContext,
			req: pb.EntityRequest{
				Data: &types.Entity{
					ID:     proto.String("entity1"),
					Secret: proto.String("secret1"),
				},
			},
			readonly: false,
			wantErr:  nil,
		},
		{
			// Fails: bad token
			ctx: InvalidAuthContext,
			req: pb.EntityRequest{
				Data: &types.Entity{
					ID:     proto.String("entity1"),
					Secret: proto.String("secret1"),
				},
			},
			readonly: false,
			wantErr:  ErrUnauthenticated,
		},
		{
			// Fails: bad permissions
			ctx: UnprivilegedContext,
			req: pb.EntityRequest{
				Data: &types.Entity{
					ID:     proto.String("entity1"),
					Secret: proto.String("secret1"),
				},
			},
			readonly: false,
			wantErr:  ErrRequestorUnqualified,
		},
		{
			// Fails: manipulation error
			ctx: PrivilegedContext,
			req: pb.EntityRequest{
				Data: &types.Entity{
					ID:     proto.String("entity1"),
					Secret: proto.String("return-error"),
				},
			},
			readonly: false,
			wantErr:  ErrInternal,
		},
	}

	for i, c := range cases {
		s := newServer(t)
		initTree(t, s)
		s.readonly = c.readonly
		if _, err := s.AuthChangeSecret(c.ctx, &c.req); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}
