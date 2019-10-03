package rpc2

import (
	"context"
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/internal/token/null"

	pb "github.com/NetAuth/Protocol/v2"
	types "github.com/NetAuth/Protocol"
)

func TestGroupCreate(t *testing.T) {
	cases := []struct{
		req pb.GroupRequest
		wantErr error
		readonly bool
	}{
		{
			// Works, valid and authorized request
			req: pb.GroupRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Group: &types.Group{
					Name: proto.String("test1"),
				},
			},
			wantErr: nil,
			readonly: false,
		},
		{
			// Fails, server is read-only
			req: pb.GroupRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Group: &types.Group{
					Name: proto.String("test1"),
				},
			},
			wantErr: ErrReadOnly,
			readonly: true,
		},
		{
			// Fails, unauthorized
			req: pb.GroupRequest{
				Group: &types.Group{
					Name: proto.String("test1"),
				},
			},
			wantErr: ErrUnauthenticated,
			readonly: false,
		},
		{
			// Fails, empty token
			req: pb.GroupRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidEmptyToken,
				},
				Group: &types.Group{
					Name: proto.String("test1"),
				},
			},
			wantErr: ErrRequestorUnqualified,
			readonly: false,
		},
		{
			// Fails, Duplicate name
			req: pb.GroupRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Group: &types.Group{
					Name: proto.String("group1"),
				},
			},
			wantErr: ErrExists,
			readonly: false,
		},
		{
			// Fails, can't be saved
			req: pb.GroupRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Group: &types.Group{
					Name: proto.String("save-error"),
				},
			},
			wantErr: ErrInternal,
			readonly: false,
		},
	}

	for i, c := range cases {
		s := newServer(t)
		initTree(t, s)
		s.readonly = c.readonly
		if _, err := s.GroupCreate(context.Background(), &c.req); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}
