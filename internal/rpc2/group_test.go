package rpc2

import (
	"context"
	"testing"

	"github.com/netauth/netauth/internal/db"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"

	types "github.com/netauth/protocol"
	pb "github.com/netauth/protocol/v2"
)

func TestGroupCreate(t *testing.T) {
	cases := []struct {
		ctx      context.Context
		req      pb.GroupRequest
		wantErr  error
		readonly bool
	}{
		{
			// Works, valid and authorized request
			ctx: PrivilegedContext,
			req: pb.GroupRequest{
				Group: &types.Group{
					Name: proto.String("test1"),
				},
			},
			wantErr:  nil,
			readonly: false,
		},
		{
			// Fails, server is read-only
			ctx: PrivilegedContext,
			req: pb.GroupRequest{
				Group: &types.Group{
					Name: proto.String("test1"),
				},
			},
			wantErr:  ErrReadOnly,
			readonly: true,
		},
		{
			// Fails, unauthorized
			ctx: InvalidAuthContext,
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
			ctx: UnprivilegedContext,
			req: pb.GroupRequest{
				Group: &types.Group{
					Name: proto.String("test1"),
				},
			},
			wantErr:  ErrRequestorUnqualified,
			readonly: false,
		},
		{
			// Fails, Duplicate name
			ctx: PrivilegedContext,
			req: pb.GroupRequest{
				Group: &types.Group{
					Name: proto.String("group1"),
				},
			},
			wantErr:  ErrExists,
			readonly: false,
		},
		{
			// Fails, can't be saved
			ctx: PrivilegedContext,
			req: pb.GroupRequest{
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
		initTree(t, s.Manager)
		s.readonly = c.readonly
		if _, err := s.GroupCreate(c.ctx, &c.req); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}

func TestGroupUpdate(t *testing.T) {
	cases := []struct {
		ctx      context.Context
		req      pb.GroupRequest
		wantErr  error
		readonly bool
	}{
		{
			// Works, valid and authorized request
			ctx: PrivilegedContext,
			req: pb.GroupRequest{
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
			ctx: PrivilegedContext,
			req: pb.GroupRequest{
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
			ctx: InvalidAuthContext,
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
			ctx: UnprivilegedContext,
			req: pb.GroupRequest{
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
			ctx: PrivilegedContext,
			req: pb.GroupRequest{
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
			ctx: PrivilegedContext,
			req: pb.GroupRequest{
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
		initTree(t, s.Manager)
		s.readonly = c.readonly
		if _, err := s.GroupUpdate(c.ctx, &c.req); err != c.wantErr {
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
		initTree(t, s.Manager)
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
		ctx      context.Context
		req      pb.KVRequest
		wantErr  error
		readonly bool
		wantRes  string
	}{
		{
			// Works, is an authorized write
			ctx: PrivilegedContext,
			req: pb.KVRequest{
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
			ctx: UnprivilegedContext,
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
			ctx: PrivilegedContext,
			req: pb.KVRequest{
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
			ctx: InvalidAuthContext,
			req: pb.KVRequest{
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
			ctx: UnprivilegedContext,
			req: pb.KVRequest{
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
			ctx: PrivilegedContext,
			req: pb.KVRequest{
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
			ctx: PrivilegedContext,
			req: pb.KVRequest{
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
			ctx: PrivilegedContext,
			req: pb.KVRequest{
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
		initTree(t, s.Manager)
		s.CreateGroup(c.ctx, "load-error", "", "", -1)
		s.readonly = c.readonly
		_, err := s.GroupUM(c.ctx, &c.req)
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

func TestGroupKVGet(t *testing.T) {
	cases := []struct {
		req     *pb.KV2Request
		wantErr error
		wantRes *pb.ListOfKVData
	}{
		{
			req: &pb.KV2Request{
				Target: proto.String("group1"),
				Data:   &types.KVData{Key: proto.String("key1")},
			},
			wantErr: nil,
			wantRes: &pb.ListOfKVData{KVData: []*types.KVData{{
				Key:    proto.String("key1"),
				Values: []*types.KVValue{{Value: proto.String("value1")}}}},
			},
		},
		{
			req:     &pb.KV2Request{Target: proto.String("group1")},
			wantErr: ErrDoesNotExist,
			wantRes: &pb.ListOfKVData{},
		},
		{
			req:     &pb.KV2Request{Target: proto.String("unknown")},
			wantErr: ErrDoesNotExist,
			wantRes: &pb.ListOfKVData{},
		},
		{
			req:     &pb.KV2Request{Target: proto.String("load-error")},
			wantErr: ErrInternal,
			wantRes: &pb.ListOfKVData{},
		},
	}

	for i, c := range cases {
		s := newServer(t)
		initTree(t, s.Manager)

		res, err := s.GroupKVGet(UnprivilegedContext, c.req)
		assert.Equalf(t, c.wantErr, err, "Test case %d", i)
		assert.Equalf(t, c.wantRes, res, "Test case %d", i)
	}
}

func TestGroupKVAdd(t *testing.T) {
	cases := []struct {
		ro      bool
		ctx     context.Context
		req     *pb.KV2Request
		wantErr error
	}{
		{
			ro:  false,
			ctx: PrivilegedContext,
			req: &pb.KV2Request{
				Target: proto.String("group1"),
				Data: &types.KVData{
					Key: proto.String("key2"),
				},
			},
			wantErr: nil,
		},
		{
			ro:      false,
			ctx:     UnprivilegedContext,
			req:     &pb.KV2Request{Target: proto.String("group1")},
			wantErr: ErrRequestorUnqualified,
		},
		{
			ro:      false,
			ctx:     InvalidAuthContext,
			req:     &pb.KV2Request{Target: proto.String("group1")},
			wantErr: ErrUnauthenticated,
		},
		{
			ro:      false,
			ctx:     PrivilegedContext,
			req:     &pb.KV2Request{Target: proto.String("unknown")},
			wantErr: ErrDoesNotExist,
		},
		{
			ro:      false,
			ctx:     PrivilegedContext,
			req:     &pb.KV2Request{Target: proto.String("load-error")},
			wantErr: ErrInternal,
		},
		{
			ro:      true,
			ctx:     PrivilegedContext,
			req:     &pb.KV2Request{Target: proto.String("group1")},
			wantErr: ErrReadOnly,
		},
	}

	for i, c := range cases {
		s := newServer(t)
		initTree(t, s.Manager)
		s.readonly = c.ro

		_, err := s.GroupKVAdd(c.ctx, c.req)
		assert.Equalf(t, c.wantErr, err, "Test Case %d", i)
	}
}

func TestGroupKVDel(t *testing.T) {
	cases := []struct {
		ro      bool
		ctx     context.Context
		req     *pb.KV2Request
		wantErr error
	}{
		{
			ro:  false,
			ctx: PrivilegedContext,
			req: &pb.KV2Request{
				Target: proto.String("group1"),
				Data: &types.KVData{
					Key: proto.String("key1"),
				},
			},
			wantErr: nil,
		},
		{
			ro:      false,
			ctx:     UnprivilegedContext,
			req:     &pb.KV2Request{Target: proto.String("group1")},
			wantErr: ErrRequestorUnqualified,
		},
		{
			ro:      false,
			ctx:     InvalidAuthContext,
			req:     &pb.KV2Request{Target: proto.String("group1")},
			wantErr: ErrUnauthenticated,
		},
		{
			ro:      false,
			ctx:     PrivilegedContext,
			req:     &pb.KV2Request{Target: proto.String("unknown")},
			wantErr: ErrDoesNotExist,
		},
		{
			ro:      false,
			ctx:     PrivilegedContext,
			req:     &pb.KV2Request{Target: proto.String("load-error")},
			wantErr: ErrInternal,
		},
		{
			ro:      true,
			ctx:     PrivilegedContext,
			req:     &pb.KV2Request{Target: proto.String("group1")},
			wantErr: ErrReadOnly,
		},
	}

	for i, c := range cases {
		s := newServer(t)
		initTree(t, s.Manager)
		s.Manager.GroupKVAdd(c.ctx, "group1", []*types.KVData{
			{
				Key: proto.String("key1"),
				Values: []*types.KVValue{
					{Value: proto.String("value1")},
				},
			},
		})
		s.readonly = c.ro

		_, err := s.GroupKVDel(c.ctx, c.req)
		assert.Equalf(t, c.wantErr, err, "Test Case %d", i)
	}
}

func TestGroupKVReplace(t *testing.T) {
	cases := []struct {
		ro      bool
		ctx     context.Context
		req     *pb.KV2Request
		wantErr error
	}{
		{
			ro:  false,
			ctx: PrivilegedContext,
			req: &pb.KV2Request{
				Target: proto.String("group1"),
				Data: &types.KVData{
					Key: proto.String("key1"),
				},
			},
			wantErr: nil,
		},
		{
			ro:      false,
			ctx:     UnprivilegedContext,
			req:     &pb.KV2Request{Target: proto.String("group1")},
			wantErr: ErrRequestorUnqualified,
		},
		{
			ro:      false,
			ctx:     InvalidAuthContext,
			req:     &pb.KV2Request{Target: proto.String("group1")},
			wantErr: ErrUnauthenticated,
		},
		{
			ro:      false,
			ctx:     PrivilegedContext,
			req:     &pb.KV2Request{Target: proto.String("unknown")},
			wantErr: ErrDoesNotExist,
		},
		{
			ro:      false,
			ctx:     PrivilegedContext,
			req:     &pb.KV2Request{Target: proto.String("load-error")},
			wantErr: ErrInternal,
		},
		{
			ro:      true,
			ctx:     PrivilegedContext,
			req:     &pb.KV2Request{Target: proto.String("group1")},
			wantErr: ErrReadOnly,
		},
	}

	for i, c := range cases {
		s := newServer(t)
		initTree(t, s.Manager)
		s.Manager.GroupKVAdd(c.ctx, "group1", []*types.KVData{
			{
				Key: proto.String("key1"),
				Values: []*types.KVValue{
					{Value: proto.String("value1")},
				},
			},
		})
		s.readonly = c.ro

		_, err := s.GroupKVReplace(c.ctx, c.req)
		assert.Equalf(t, c.wantErr, err, "Test Case %d", i)
	}
}

func TestGroupUpdateRules(t *testing.T) {
	cases := []struct {
		ctx      context.Context
		req      pb.GroupRulesRequest
		readonly bool
		wantErr  error
	}{
		{
			// Works
			ctx: PrivilegedContext,
			req: pb.GroupRulesRequest{
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
			ctx: InvalidAuthContext,
			req: pb.GroupRulesRequest{
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
			ctx: UnprivilegedContext,
			req: pb.GroupRulesRequest{
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
			ctx: PrivilegedContext,
			req: pb.GroupRulesRequest{
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
			ctx: PrivilegedContext,
			req: pb.GroupRulesRequest{
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
			ctx: PrivilegedContext,
			req: pb.GroupRulesRequest{
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
		initTree(t, s.Manager)
		s.readonly = c.readonly
		if _, err := s.GroupUpdateRules(c.ctx, &c.req); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}

func TestGroupAddMember(t *testing.T) {
	cases := []struct {
		ctx      context.Context
		req      pb.EntityRequest
		wantErr  error
		readonly bool
	}{
		{
			// Works
			ctx: PrivilegedContext,
			req: pb.EntityRequest{
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
			ctx: PrivilegedContext,
			req: pb.EntityRequest{
				Entity: &types.Entity{
					ID: proto.String("entity1"),
				},
			},
			wantErr:  nil,
			readonly: false,
		},
		{
			// Fails, readonly
			ctx: PrivilegedContext,
			req: pb.EntityRequest{
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
			ctx: InvalidAuthContext,
			req: pb.EntityRequest{
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
			ctx: UnprivilegedContext,
			req: pb.EntityRequest{
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
			ctx: PrivilegedContext,
			req: pb.EntityRequest{
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
		initTree(t, s.Manager)
		s.readonly = c.readonly
		if _, err := s.GroupAddMember(c.ctx, &c.req); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}

func TestGroupDelMember(t *testing.T) {
	cases := []struct {
		ctx      context.Context
		req      pb.EntityRequest
		wantErr  error
		readonly bool
	}{
		{
			// Works
			ctx: PrivilegedContext,
			req: pb.EntityRequest{
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
			ctx: PrivilegedContext,
			req: pb.EntityRequest{
				Entity: &types.Entity{
					ID: proto.String("entity1"),
				},
			},
			wantErr:  nil,
			readonly: false,
		},
		{
			// Fails, readonly
			ctx: PrivilegedContext,
			req: pb.EntityRequest{
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
			ctx: InvalidAuthContext,
			req: pb.EntityRequest{
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
			ctx: UnprivilegedContext,
			req: pb.EntityRequest{
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
			ctx: PrivilegedContext,
			req: pb.EntityRequest{
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
		initTree(t, s.Manager)
		s.readonly = c.readonly
		if _, err := s.GroupDelMember(c.ctx, &c.req); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}

func TestGroupDestroy(t *testing.T) {
	cases := []struct {
		ctx      context.Context
		req      pb.GroupRequest
		wantErr  error
		readonly bool
	}{
		{
			// Works, is authorized
			ctx: PrivilegedContext,
			req: pb.GroupRequest{
				Group: &types.Group{
					Name: proto.String("group1"),
				},
			},
			wantErr:  nil,
			readonly: false,
		},
		{
			// Fails, read only
			ctx: PrivilegedContext,
			req: pb.GroupRequest{
				Group: &types.Group{
					Name: proto.String("group1"),
				},
			},
			wantErr:  ErrReadOnly,
			readonly: true,
		},
		{
			// Fails, bad token
			ctx: InvalidAuthContext,
			req: pb.GroupRequest{
				Group: &types.Group{
					Name: proto.String("group1"),
				},
			},
			wantErr:  ErrUnauthenticated,
			readonly: false,
		},
		{
			// Fails, empty token
			ctx: UnprivilegedContext,
			req: pb.GroupRequest{
				Group: &types.Group{
					Name: proto.String("group1"),
				},
			},
			wantErr:  ErrRequestorUnqualified,
			readonly: false,
		},
		{
			// Fails, unknown group
			ctx: PrivilegedContext,
			req: pb.GroupRequest{
				Group: &types.Group{
					Name: proto.String("does-not-exist"),
				},
			},
			wantErr:  ErrDoesNotExist,
			readonly: false,
		},
		{
			// Fails, load-error
			ctx: PrivilegedContext,
			req: pb.GroupRequest{
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
		initTree(t, s.Manager)
		s.readonly = c.readonly
		if _, err := s.GroupDestroy(c.ctx, &c.req); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}

func TestGroupMembers(t *testing.T) {
	cases := []struct {
		group      string
		wantErr    error
		wantMember string
	}{
		{
			group:      "group1",
			wantErr:    nil,
			wantMember: "entity1",
		},
		{
			group:      "does-not-exist",
			wantErr:    nil,
			wantMember: "",
		},
		{
			group:      "load-error",
			wantErr:    nil,
			wantMember: "",
		},
	}

	for i, c := range cases {
		s := newServer(t)
		initTree(t, s.Manager)

		req := pb.GroupRequest{
			Group: &types.Group{
				Name: &c.group,
			},
		}

		res, err := s.GroupMembers(context.Background(), &req)
		if err != c.wantErr {
			t.Errorf("%d (%s): Got %v; Want %v", i, c.group, err, c.wantErr)
		}
		if err != nil || len(res.GetEntities()) < 1 {
			continue
		}
		if res.GetEntities()[0].GetID() != c.wantMember {
			t.Errorf("%d: Got %v; Wanted %v as member", i, res.GetEntities(), c.wantMember)
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
		s, d, _ := newServerWithRefs(t)
		initTree(t, s.Manager)
		d.(*db.DB).IndexGroup(&types.Group{Name: proto.String("load-error")})
		if _, err := s.GroupSearch(context.Background(), &pb.SearchRequest{Expression: &c.expr}); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}
