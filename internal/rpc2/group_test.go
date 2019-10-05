package rpc2

import (
	"context"
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/internal/token/null"

	types "github.com/NetAuth/Protocol"
	pb "github.com/NetAuth/Protocol/v2"
)

func TestGroupCreate(t *testing.T) {
	cases := []struct {
		req      pb.GroupRequest
		wantErr  error
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
			wantErr:  nil,
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
			wantErr:  ErrReadOnly,
			readonly: true,
		},
		{
			// Fails, unauthorized
			req: pb.GroupRequest{
				Group: &types.Group{
					Name: proto.String("test1"),
				},
			},
			wantErr:  ErrUnauthenticated,
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
			wantErr:  ErrRequestorUnqualified,
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
			wantErr:  ErrExists,
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
			wantErr:  ErrInternal,
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

func TestGroupUpdate(t *testing.T) {
	cases := []struct {
		req      pb.GroupRequest
		wantErr  error
		readonly bool
	}{
		{
			// Works, valid and authorized request
			req: pb.GroupRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Group: &types.Group{
					Name:        proto.String("group1"),
					DisplayName: proto.String("First Group"),
				},
			},
			wantErr:  nil,
			readonly: false,
		},
		{
			// Fails, server is read-only
			req: pb.GroupRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Group: &types.Group{
					Name:        proto.String("group1"),
					DisplayName: proto.String("First Group"),
				},
			},
			wantErr:  ErrReadOnly,
			readonly: true,
		},
		{
			// Fails, unauthorized
			req: pb.GroupRequest{
				Group: &types.Group{
					Name:        proto.String("group1"),
					DisplayName: proto.String("First Group"),
				},
			},
			wantErr:  ErrUnauthenticated,
			readonly: false,
		},
		{
			// Fails, empty token
			req: pb.GroupRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidEmptyToken,
				},
				Group: &types.Group{
					Name:        proto.String("group1"),
					DisplayName: proto.String("First Group"),
				},
			},
			wantErr:  ErrRequestorUnqualified,
			readonly: false,
		},
		{
			// Fails, unknown group
			req: pb.GroupRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Group: &types.Group{
					Name:        proto.String("does-not-exit"),
					DisplayName: proto.String("First Group"),
				},
			},
			wantErr:  ErrDoesNotExist,
			readonly: false,
		},
		{
			// Fails, can't be loaded
			req: pb.GroupRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Group: &types.Group{
					Name:        proto.String("load-error"),
					DisplayName: proto.String("First Group"),
				},
			},
			wantErr:  ErrInternal,
			readonly: false,
		},
	}

	for i, c := range cases {
		s := newServer(t)
		initTree(t, s)
		s.readonly = c.readonly
		if _, err := s.GroupUpdate(context.Background(), &c.req); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}

func TestGroupInfo(t *testing.T) {
	cases := []struct {
		req     pb.GroupRequest
		wantErr error
		wantLen int
	}{
		{
			req: pb.GroupRequest{
				Group: &types.Group{
					Name: proto.String("group1"),
				},
			},
			wantErr: nil,
			wantLen: 1,
		},
		{
			req: pb.GroupRequest{
				Group: &types.Group{
					Name: proto.String("does-not-exist"),
				},
			},
			wantErr: ErrDoesNotExist,
			wantLen: 0,
		},
		{
			req: pb.GroupRequest{
				Group: &types.Group{
					Name: proto.String("load-error"),
				},
			},
			wantErr: ErrInternal,
			wantLen: 0,
		},
	}

	for i, c := range cases {
		s := newServer(t)
		initTree(t, s)
		resp, err := s.GroupInfo(context.Background(), &c.req)
		if err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
		if len(resp.GetGroups()) != c.wantLen {
			t.Errorf("%d: Got %d; Want %d", i, len(resp.GetGroups()), c.wantLen)
		}
	}
}

func TestGroupUM(t *testing.T) {
	cases := []struct {
		req      pb.KVRequest
		wantErr  error
		readonly bool
		wantRes  string
	}{
		{
			// Works, is an authorized write
			req: pb.KVRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Target: proto.String("group1"),
				Action: pb.Action_UPSERT.Enum(),
				Key:    proto.String("key1"),
				Value:  proto.String("value1"),
			},
			wantErr:  nil,
			readonly: false,
			wantRes:  "key1:value1",
		},
		{
			// Works, is a read-only query
			req: pb.KVRequest{
				Target: proto.String("group1"),
				Action: pb.Action_READ.Enum(),
				Key:    proto.String("key1"),
			},
			wantErr:  nil,
			readonly: false,
			wantRes:  "",
		},
		{
			// Fails, Server is read-only
			req: pb.KVRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Target: proto.String("group1"),
				Action: pb.Action_UPSERT.Enum(),
				Key:    proto.String("key1"),
				Value:  proto.String("value1"),
			},
			wantErr:  ErrReadOnly,
			readonly: true,
			wantRes:  "",
		},
		{
			// Fails, Token is invalid
			req: pb.KVRequest{
				Auth: &pb.AuthData{
					Token: &null.InvalidToken,
				},
				Target: proto.String("group1"),
				Action: pb.Action_UPSERT.Enum(),
				Key:    proto.String("key1"),
				Value:  proto.String("value1"),
			},
			wantErr:  ErrUnauthenticated,
			readonly: false,
			wantRes:  "",
		},
		{
			// Fails, Token has no capability
			req: pb.KVRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidEmptyToken,
				},
				Target: proto.String("group1"),
				Action: pb.Action_UPSERT.Enum(),
				Key:    proto.String("key1"),
				Value:  proto.String("value1"),
			},
			wantErr:  ErrRequestorUnqualified,
			readonly: false,
			wantRes:  "",
		},
		{
			// Fails, group doesn't exist
			req: pb.KVRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Target: proto.String("does-not-exist"),
				Action: pb.Action_UPSERT.Enum(),
				Key:    proto.String("key1"),
				Value:  proto.String("value1"),
			},
			wantErr:  ErrDoesNotExist,
			readonly: false,
			wantRes:  "",
		},
		{
			// Fails, failure during load
			req: pb.KVRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Target: proto.String("load-error"),
				Action: pb.Action_UPSERT.Enum(),
				Key:    proto.String("key1"),
				Value:  proto.String("value1"),
			},
			wantErr:  ErrInternal,
			readonly: false,
			wantRes:  "",
		},
		{
			// Fails, bad request
			req: pb.KVRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Target: proto.String("group1"),
				Action: pb.Action_ADD.Enum(),
				Key:    proto.String("key1"),
				Value:  proto.String("value1"),
			},
			wantErr:  ErrMalformedRequest,
			readonly: false,
			wantRes:  "",
		},
	}

	for i, c := range cases {
		s := newServer(t)
		initTree(t, s)
		s.CreateGroup("load-error", "", "", -1)
		s.readonly = c.readonly
		_, err := s.GroupUM(context.Background(), &c.req)
		if err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
		c.req.Action = pb.Action_READ.Enum()
		res, err := s.GroupUM(context.Background(), &c.req)
		if err != nil && err != c.wantErr {
			t.Fatalf("%d: Error on readback: %v", i, err)
		}

		if len(res.GetStrings()) != 0 && res.GetStrings()[0] != c.wantRes {
			t.Errorf("%d: Got '%s'; Want '%s'", i, res, c.wantRes)
		}
	}
}

func TestGroupUpdateRules(t *testing.T) {
	cases := []struct {
		req      pb.GroupRulesRequest
		readonly bool
		wantErr  error
	}{
		{
			// Works
			req: pb.GroupRulesRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Group: &types.Group{
					Name: proto.String("group1"),
				},
				Target: &types.Group{
					Name: proto.String("group2"),
				},
				RuleAction: pb.RuleAction_INCLUDE.Enum(),
			},
			readonly: false,
			wantErr:  nil,
		},
		{
			// Fails, bad token
			req: pb.GroupRulesRequest{
				Auth: &pb.AuthData{
					Token: &null.InvalidToken,
				},
				Group: &types.Group{
					Name: proto.String("group1"),
				},
				Target: &types.Group{
					Name: proto.String("group2"),
				},
				RuleAction: pb.RuleAction_INCLUDE.Enum(),
			},
			readonly: false,
			wantErr:  ErrUnauthenticated,
		},
		{
			// Fails, empty
			req: pb.GroupRulesRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidEmptyToken,
				},
				Group: &types.Group{
					Name: proto.String("group1"),
				},
				Target: &types.Group{
					Name: proto.String("group2"),
				},
				RuleAction: pb.RuleAction_INCLUDE.Enum(),
			},
			readonly: false,
			wantErr:  ErrRequestorUnqualified,
		},
		{
			// Fails, read-only
			req: pb.GroupRulesRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Group: &types.Group{
					Name: proto.String("group1"),
				},
				Target: &types.Group{
					Name: proto.String("group2"),
				},
				RuleAction: pb.RuleAction_INCLUDE.Enum(),
			},
			readonly: true,
			wantErr:  ErrReadOnly,
		},
		{
			// Fails, does not exist
			req: pb.GroupRulesRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Group: &types.Group{
					Name: proto.String("group1"),
				},
				Target: &types.Group{
					Name: proto.String("does-not-exist"),
				},
				RuleAction: pb.RuleAction_INCLUDE.Enum(),
			},
			readonly: false,
			wantErr:  ErrDoesNotExist,
		},
		{
			// Fails, load-error
			req: pb.GroupRulesRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Group: &types.Group{
					Name: proto.String("group1"),
				},
				Target: &types.Group{
					Name: proto.String("load-error"),
				},
				RuleAction: pb.RuleAction_INCLUDE.Enum(),
			},
			readonly: false,
			wantErr:  ErrInternal,
		},
	}

	for i, c := range cases {
		s := newServer(t)
		initTree(t, s)
		s.readonly = c.readonly
		if _, err := s.GroupUpdateRules(context.Background(), &c.req); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}

func TestGroupDelMember(t *testing.T) {
	cases := []struct {
		req      pb.EntityRequest
		wantErr  error
		readonly bool
	}{
		{
			// Works
			req: pb.EntityRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Entity: &types.Entity{
					ID: proto.String("entity1"),
					Meta: &types.EntityMeta{
						Groups: []string{
							"group1",
						},
					},
				},
			},
			wantErr:  nil,
			readonly: false,
		},
		{
			// Works, no groups
			req: pb.EntityRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Entity: &types.Entity{
					ID: proto.String("entity1"),
				},
			},
			wantErr:  nil,
			readonly: false,
		},
		{
			// Fails, readonly
			req: pb.EntityRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Entity: &types.Entity{
					ID: proto.String("entity1"),
					Meta: &types.EntityMeta{
						Groups: []string{
							"group1",
						},
					},
				},
			},
			wantErr:  ErrReadOnly,
			readonly: true,
		},
		{
			// Fails, bad token
			req: pb.EntityRequest{
				Auth: &pb.AuthData{
					Token: &null.InvalidToken,
				},
				Entity: &types.Entity{
					ID: proto.String("entity1"),
					Meta: &types.EntityMeta{
						Groups: []string{
							"group1",
						},
					},
				},
			},
			wantErr:  ErrUnauthenticated,
			readonly: false,
		},
		{
			// Fails, empty token
			req: pb.EntityRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidEmptyToken,
				},
				Entity: &types.Entity{
					ID: proto.String("entity1"),
					Meta: &types.EntityMeta{
						Groups: []string{
							"group1",
						},
					},
				},
			},
			wantErr:  ErrRequestorUnqualified,
			readonly: false,
		},
		{
			// Fails, entity can't be loaded
			req: pb.EntityRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Entity: &types.Entity{
					ID: proto.String("load-error"),
					Meta: &types.EntityMeta{
						Groups: []string{
							"group1",
						},
					},
				},
			},
			wantErr:  ErrInternal,
			readonly: false,
		},
	}
	for i, c := range cases {
		s := newServer(t)
		initTree(t, s)
		s.readonly = c.readonly
		if _, err := s.GroupDelMember(context.Background(), &c.req); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}

func TestGroupAddMember(t *testing.T) {
	cases := []struct {
		req      pb.EntityRequest
		wantErr  error
		readonly bool
	}{
		{
			// Works
			req: pb.EntityRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Entity: &types.Entity{
					ID: proto.String("entity1"),
					Meta: &types.EntityMeta{
						Groups: []string{
							"group1",
						},
					},
				},
			},
			wantErr:  nil,
			readonly: false,
		},
		{
			// Works, no groups
			req: pb.EntityRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Entity: &types.Entity{
					ID: proto.String("entity1"),
				},
			},
			wantErr:  nil,
			readonly: false,
		},
		{
			// Fails, readonly
			req: pb.EntityRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Entity: &types.Entity{
					ID: proto.String("entity1"),
					Meta: &types.EntityMeta{
						Groups: []string{
							"group1",
						},
					},
				},
			},
			wantErr:  ErrReadOnly,
			readonly: true,
		},
		{
			// Fails, bad token
			req: pb.EntityRequest{
				Auth: &pb.AuthData{
					Token: &null.InvalidToken,
				},
				Entity: &types.Entity{
					ID: proto.String("entity1"),
					Meta: &types.EntityMeta{
						Groups: []string{
							"group1",
						},
					},
				},
			},
			wantErr:  ErrUnauthenticated,
			readonly: false,
		},
		{
			// Fails, empty token
			req: pb.EntityRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidEmptyToken,
				},
				Entity: &types.Entity{
					ID: proto.String("entity1"),
					Meta: &types.EntityMeta{
						Groups: []string{
							"group1",
						},
					},
				},
			},
			wantErr:  ErrRequestorUnqualified,
			readonly: false,
		},
		{
			// Fails, entity can't be loaded
			req: pb.EntityRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Entity: &types.Entity{
					ID: proto.String("load-error"),
					Meta: &types.EntityMeta{
						Groups: []string{
							"group1",
						},
					},
				},
			},
			wantErr:  ErrInternal,
			readonly: false,
		},
	}
	for i, c := range cases {
		s := newServer(t)
		initTree(t, s)
		s.readonly = c.readonly
		if _, err := s.GroupAddMember(context.Background(), &c.req); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}

func TestGroupDestroy(t *testing.T) {
	cases := []struct {
		req      pb.GroupRequest
		wantErr  error
		readonly bool
	}{
		{
			// Works, is authorized
			req: pb.GroupRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Group: &types.Group{
					Name: proto.String("group1"),
				},
			},
			wantErr:  nil,
			readonly: false,
		},
		{
			// Fails, read only
			req: pb.GroupRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Group: &types.Group{
					Name: proto.String("group1"),
				},
			},
			wantErr:  ErrReadOnly,
			readonly: true,
		},
		{
			// Fails, bad token
			req: pb.GroupRequest{
				Auth: &pb.AuthData{
					Token: &null.InvalidToken,
				},
				Group: &types.Group{
					Name: proto.String("group1"),
				},
			},
			wantErr:  ErrUnauthenticated,
			readonly: false,
		},
		{
			// Fails, empty token
			req: pb.GroupRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidEmptyToken,
				},
				Group: &types.Group{
					Name: proto.String("group1"),
				},
			},
			wantErr:  ErrRequestorUnqualified,
			readonly: false,
		},
		{
			// Fails, unknown group
			req: pb.GroupRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Group: &types.Group{
					Name: proto.String("does-not-exist"),
				},
			},
			wantErr:  ErrDoesNotExist,
			readonly: false,
		},
		{
			// Fails, load-error
			req: pb.GroupRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Group: &types.Group{
					Name: proto.String("load-error"),
				},
			},
			wantErr:  ErrInternal,
			readonly: false,
		},
	}

	for i, c := range cases {
		s := newServer(t)
		initTree(t, s)
		s.readonly = c.readonly
		if _, err := s.GroupDestroy(context.Background(), &c.req); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}

func TestGroupSearch(t *testing.T) {
	cases := []struct {
		expr    string
		wantErr error
	}{
		{
			expr:    "group1",
			wantErr: nil,
		},
		{
			expr:    "*",
			wantErr: ErrInternal,
		},
	}

	for i, c := range cases {
		s := newServer(t)
		initTree(t, s)
		s.CreateGroup("load-error", "", "", -1)
		if _, err := s.GroupSearch(context.Background(), &pb.SearchRequest{Expression: &c.expr}); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}
