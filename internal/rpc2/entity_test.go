package rpc2

import (
	"context"
	"testing"

	"github.com/golang/protobuf/proto"

	types "github.com/NetAuth/Protocol"
	pb "github.com/NetAuth/Protocol/v2"
)

func TestEntityCreate(t *testing.T) {
	cases := []struct {
		req      pb.EntityRequest
		wantErr  error
		readonly bool
	}{
		{
			// Works, entity is created.
			req: pb.EntityRequest{
				Auth: ValidAuthData,
				Entity: &types.Entity{
					ID: proto.String("test1"),
				},
			},
			wantErr:  nil,
			readonly: false,
		},
		{
			// Fails, server is in read-only mode
			req: pb.EntityRequest{
				Auth: ValidAuthData,
				Entity: &types.Entity{
					ID: proto.String("test1"),
				},
			},
			wantErr:  ErrReadOnly,
			readonly: true,
		},
		{
			// Fails, token is invalid
			req: pb.EntityRequest{
				Auth: InvalidAuthData,
				Entity: &types.Entity{
					ID: proto.String("test1"),
				},
			},
			wantErr:  ErrUnauthenticated,
			readonly: false,
		},
		{
			// Fails, token lacks capabilities
			req: pb.EntityRequest{
				Auth: EmptyAuthData,
				Entity: &types.Entity{
					ID: proto.String("test1"),
				},
			},
			wantErr:  ErrRequestorUnqualified,
			readonly: false,
		},
		{
			// Fails, duplicate resource
			req: pb.EntityRequest{
				Auth: ValidAuthData,
				Entity: &types.Entity{
					// This gets created by
					// initTree which fills in the
					// tree for testing purposes.
					ID: proto.String("entity1"),
				},
			},
			wantErr:  ErrExists,
			readonly: false,
		},
		{
			// Fails, internal write error
			req: pb.EntityRequest{
				Auth: ValidAuthData,
				Entity: &types.Entity{
					ID: proto.String("save-error"),
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
		if _, err := s.EntityCreate(context.Background(), &c.req); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}

func TestEntityUpdate(t *testing.T) {
	cases := []struct {
		req      pb.EntityRequest
		readonly bool
		wantErr  error
	}{
		{
			// Works, will change the metadata DisplayName
			req: pb.EntityRequest{
				Auth: ValidAuthData,
				Data: &types.Entity{
					ID: proto.String("entity1"),
					Meta: &types.EntityMeta{
						DisplayName: proto.String("First Entity"),
					},
				},
			},
			readonly: false,
			wantErr:  nil,
		},
		{
			// Fails, server is in read-only mode
			req: pb.EntityRequest{
				Auth: ValidAuthData,
				Data: &types.Entity{
					ID: proto.String("entity1"),
					Meta: &types.EntityMeta{
						DisplayName: proto.String("First Entity"),
					},
				},
			},
			readonly: true,
			wantErr:  ErrReadOnly,
		},
		{
			// Fails, token is invalid
			req: pb.EntityRequest{
				Auth: InvalidAuthData,
				Data: &types.Entity{
					ID: proto.String("entity1"),
					Meta: &types.EntityMeta{
						DisplayName: proto.String("First Entity"),
					},
				},
			},
			readonly: false,
			wantErr:  ErrUnauthenticated,
		},
		{
			// Fails, token has no capabilities
			req: pb.EntityRequest{
				Auth: EmptyAuthData,
				Data: &types.Entity{
					ID: proto.String("entity1"),
					Meta: &types.EntityMeta{
						DisplayName: proto.String("First Entity"),
					},
				},
			},
			readonly: false,
			wantErr:  ErrRequestorUnqualified,
		},
		{
			// Fails, entity does not exist
			req: pb.EntityRequest{
				Auth: ValidAuthData,
				Data: &types.Entity{
					ID: proto.String("does-not-exist"),
					Meta: &types.EntityMeta{
						DisplayName: proto.String("First Entity"),
					},
				},
			},
			readonly: false,
			wantErr:  ErrDoesNotExist,
		},
		{
			// Fails, db write failure
			req: pb.EntityRequest{
				Auth: ValidAuthData,
				Data: &types.Entity{
					ID: proto.String("load-error"),
					Meta: &types.EntityMeta{
						DisplayName: proto.String("First Entity"),
					},
				},
			},
			readonly: false,
			wantErr:  ErrInternal,
		},
	}

	for i, c := range cases {
		s := newServer(t)
		s.readonly = c.readonly
		initTree(t, s)
		if _, err := s.EntityUpdate(context.Background(), &c.req); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}

func TestEntityInfo(t *testing.T) {
	cases := []struct {
		req     pb.EntityRequest
		wantErr error
		wantLen int
	}{
		{
			// Works
			req: pb.EntityRequest{
				Entity: &types.Entity{
					ID: proto.String("entity1"),
				},
			},
			wantErr: nil,
			wantLen: 1,
		},
		{
			// Fails, does not exist
			req: pb.EntityRequest{
				Entity: &types.Entity{
					ID: proto.String("does-not-exist"),
				},
			},
			wantErr: ErrDoesNotExist,
			wantLen: 0,
		},
		{
			// Fails, load-error
			req: pb.EntityRequest{
				Entity: &types.Entity{
					ID: proto.String("load-error"),
				},
			},
			wantErr: ErrInternal,
			wantLen: 0,
		},
	}

	for i, c := range cases {
		s := newServer(t)
		initTree(t, s)
		if res, err := s.EntityInfo(context.Background(), &c.req); err != c.wantErr || len(res.Entities) != c.wantLen {
			t.Errorf("%d: Got %d, %v; Want %d, %v", i, len(res.Entities), err, c.wantLen, c.wantErr)
		}
	}
}

func TestEntitySearch(t *testing.T) {
	cases := []struct {
		req     pb.SearchRequest
		wantErr error
	}{
		{
			// Works, entity1 can be loaded
			req: pb.SearchRequest{
				Expression: proto.String("ID:entity1"),
			},
			wantErr: nil,
		},
		{
			// Fails, load-error is included in the set of all
			req: pb.SearchRequest{
				Expression: proto.String("*"),
			},
			wantErr: ErrInternal,
		},
	}

	for i, c := range cases {
		s := newServer(t)
		initTree(t, s)
		s.CreateEntity("load-error", -1, "")
		if _, err := s.EntitySearch(context.Background(), &c.req); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}

func TestEntityUM(t *testing.T) {
	cases := []struct {
		req      pb.KVRequest
		wantErr  error
		readonly bool
		wantRes  string
	}{
		{
			// Works, is an authorized write
			req: pb.KVRequest{
				Auth:   ValidAuthData,
				Target: proto.String("entity1"),
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
				Target: proto.String("entity1"),
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
				Auth:   ValidAuthData,
				Target: proto.String("entity1"),
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
				Auth:   InvalidAuthData,
				Target: proto.String("entity1"),
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
				Auth:   EmptyAuthData,
				Target: proto.String("entity1"),
				Action: pb.Action_UPSERT.Enum(),
				Key:    proto.String("key1"),
				Value:  proto.String("value1"),
			},
			wantErr:  ErrRequestorUnqualified,
			readonly: false,
			wantRes:  "",
		},
		{
			// Fails, entity doesn't exist
			req: pb.KVRequest{
				Auth:   ValidAuthData,
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
				Auth:   ValidAuthData,
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
				Auth:   ValidAuthData,
				Target: proto.String("entity1"),
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
		s.CreateEntity("load-error", -1, "")
		s.readonly = c.readonly
		_, err := s.EntityUM(context.Background(), &c.req)
		if err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
		c.req.Action = pb.Action_READ.Enum()
		res, err := s.EntityUM(context.Background(), &c.req)
		if err != nil && err != c.wantErr {
			t.Fatalf("%d: Error on readback: %v", i, err)
		}

		if len(res.GetStrings()) != 0 && res.GetStrings()[0] != c.wantRes {
			t.Errorf("%d: Got '%s'; Want '%s'", i, res, c.wantRes)
		}
	}
}

func TestEntityKeys(t *testing.T) {
	cases := []struct {
		req      pb.KVRequest
		wantErr  error
		readonly bool
		wantRes  string
	}{
		{
			// Works, is an authorized write
			req: pb.KVRequest{
				Auth:   ValidAuthData,
				Target: proto.String("entity1"),
				Action: pb.Action_ADD.Enum(),
				Key:    proto.String("ssh"),
				Value:  proto.String("key1"),
			},
			wantErr:  nil,
			readonly: false,
			wantRes:  "SSH:key1",
		},
		{
			// Works, self-change is allowed
			req: pb.KVRequest{
				Auth:   EmptyAuthData,
				Target: proto.String("valid"),
				Action: pb.Action_ADD.Enum(),
				Key:    proto.String("ssh"),
				Value:  proto.String("key1"),
			},
			wantErr:  nil,
			readonly: false,
			wantRes:  "SSH:key1",
		},
		{
			// Works, is a read-only query
			req: pb.KVRequest{
				Target: proto.String("entity1"),
				Action: pb.Action_READ.Enum(),
				Key:    proto.String("ssh"),
			},
			wantErr:  nil,
			readonly: false,
			wantRes:  "",
		},
		{
			// Fails, Server is read-only
			req: pb.KVRequest{
				Auth:   ValidAuthData,
				Target: proto.String("entity1"),
				Action: pb.Action_ADD.Enum(),
				Key:    proto.String("ssh"),
				Value:  proto.String("value1"),
			},
			wantErr:  ErrReadOnly,
			readonly: true,
			wantRes:  "",
		},
		{
			// Fails, Token is invalid
			req: pb.KVRequest{
				Auth:   InvalidAuthData,
				Target: proto.String("entity1"),
				Action: pb.Action_ADD.Enum(),
				Key:    proto.String("ssh"),
				Value:  proto.String("value1"),
			},
			wantErr:  ErrUnauthenticated,
			readonly: false,
			wantRes:  "",
		},
		{
			// Fails, Token has no capability
			req: pb.KVRequest{
				Auth:   EmptyAuthData,
				Target: proto.String("entity1"),
				Action: pb.Action_ADD.Enum(),
				Key:    proto.String("ssh"),
				Value:  proto.String("value1"),
			},
			wantErr:  ErrRequestorUnqualified,
			readonly: false,
			wantRes:  "",
		},
		{
			// Fails, entity doesn't exist
			req: pb.KVRequest{
				Auth:   ValidAuthData,
				Target: proto.String("does-not-exist"),
				Action: pb.Action_ADD.Enum(),
				Key:    proto.String("ssh"),
				Value:  proto.String("value1"),
			},
			wantErr:  ErrDoesNotExist,
			readonly: false,
			wantRes:  "",
		},
		{
			// Fails, failure during load
			req: pb.KVRequest{
				Auth:   ValidAuthData,
				Target: proto.String("load-error"),
				Action: pb.Action_ADD.Enum(),
				Key:    proto.String("ssh"),
				Value:  proto.String("value1"),
			},
			wantErr:  ErrInternal,
			readonly: false,
			wantRes:  "",
		},
		{
			// Fails, bad request
			req: pb.KVRequest{
				Auth:   ValidAuthData,
				Target: proto.String("entity1"),
				Action: pb.Action_UPSERT.Enum(),
				Key:    proto.String("ssh"),
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
		s.CreateEntity("valid", -1, "")
		s.CreateEntity("load-error", -1, "")
		s.readonly = c.readonly
		_, err := s.EntityKeys(context.Background(), &c.req)
		if err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
		c.req.Action = pb.Action_READ.Enum()
		res, err := s.EntityKeys(context.Background(), &c.req)
		if err != nil && err != c.wantErr {
			t.Fatalf("%d: Error on readback: %v", i, err)
		}

		if len(res.GetStrings()) != 0 && res.GetStrings()[0] != c.wantRes {
			t.Errorf("%d: Got '%s'; Want '%s'", i, res, c.wantRes)
		}
	}
}

func TestEntityDestroy(t *testing.T) {
	cases := []struct {
		req      pb.EntityRequest
		wantErr  error
		readonly bool
	}{
		{
			// Works, entity is created.
			req: pb.EntityRequest{
				Auth: ValidAuthData,
				Entity: &types.Entity{
					ID: proto.String("entity1"),
				},
			},
			wantErr:  nil,
			readonly: false,
		},
		{
			// Fails, server is in read-only mode
			req: pb.EntityRequest{
				Auth: ValidAuthData,
				Entity: &types.Entity{
					ID: proto.String("entity1"),
				},
			},
			wantErr:  ErrReadOnly,
			readonly: true,
		},
		{
			// Fails, token is invalid
			req: pb.EntityRequest{
				Auth: InvalidAuthData,
				Entity: &types.Entity{
					ID: proto.String("entity1"),
				},
			},
			wantErr:  ErrUnauthenticated,
			readonly: false,
		},
		{
			// Fails, token lacks capabilities
			req: pb.EntityRequest{
				Auth: EmptyAuthData,
				Entity: &types.Entity{
					ID: proto.String("entity1"),
				},
			},
			wantErr:  ErrRequestorUnqualified,
			readonly: false,
		},
		{
			// Fails, internal write error
			req: pb.EntityRequest{
				Auth: ValidAuthData,
				Entity: &types.Entity{
					ID: proto.String("load-error"),
				},
			},
			wantErr:  ErrInternal,
			readonly: false,
		},
		{
			// Fails, unknown user
			req: pb.EntityRequest{
				Auth: ValidAuthData,
				Entity: &types.Entity{
					ID: proto.String("does-not-exist"),
				},
			},
			wantErr:  ErrDoesNotExist,
			readonly: false,
		},
	}

	for i, c := range cases {
		s := newServer(t)
		initTree(t, s)
		s.readonly = c.readonly
		if _, err := s.EntityDestroy(context.Background(), &c.req); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}

func TestEntityLock(t *testing.T) {
	cases := []struct {
		req      pb.EntityRequest
		wantErr  error
		readonly bool
	}{
		{
			// Works, entity is created.
			req: pb.EntityRequest{
				Auth: ValidAuthData,
				Entity: &types.Entity{
					ID: proto.String("entity1"),
				},
			},
			wantErr:  nil,
			readonly: false,
		},
		{
			// Fails, server is in read-only mode
			req: pb.EntityRequest{
				Auth: ValidAuthData,
				Entity: &types.Entity{
					ID: proto.String("entity1"),
				},
			},
			wantErr:  ErrReadOnly,
			readonly: true,
		},
		{
			// Fails, token is invalid
			req: pb.EntityRequest{
				Auth: InvalidAuthData,
				Entity: &types.Entity{
					ID: proto.String("entity1"),
				},
			},
			wantErr:  ErrUnauthenticated,
			readonly: false,
		},
		{
			// Fails, token lacks capabilities
			req: pb.EntityRequest{
				Auth: EmptyAuthData,
				Entity: &types.Entity{
					ID: proto.String("entity1"),
				},
			},
			wantErr:  ErrRequestorUnqualified,
			readonly: false,
		},
		{
			// Fails, internal write error
			req: pb.EntityRequest{
				Auth: ValidAuthData,
				Entity: &types.Entity{
					ID: proto.String("load-error"),
				},
			},
			wantErr:  ErrInternal,
			readonly: false,
		},
		{
			// Fails, unknown user
			req: pb.EntityRequest{
				Auth: ValidAuthData,
				Entity: &types.Entity{
					ID: proto.String("does-not-exist"),
				},
			},
			wantErr:  ErrDoesNotExist,
			readonly: false,
		},
	}

	for i, c := range cases {
		s := newServer(t)
		initTree(t, s)
		s.readonly = c.readonly
		if _, err := s.EntityLock(context.Background(), &c.req); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}

func TestEntityUnlock(t *testing.T) {
	cases := []struct {
		req      pb.EntityRequest
		wantErr  error
		readonly bool
	}{
		{
			// Works, entity is created.
			req: pb.EntityRequest{
				Auth: ValidAuthData,
				Entity: &types.Entity{
					ID: proto.String("entity1"),
				},
			},
			wantErr:  nil,
			readonly: false,
		},
		{
			// Fails, server is in read-only mode
			req: pb.EntityRequest{
				Auth: ValidAuthData,
				Entity: &types.Entity{
					ID: proto.String("entity1"),
				},
			},
			wantErr:  ErrReadOnly,
			readonly: true,
		},
		{
			// Fails, token is invalid
			req: pb.EntityRequest{
				Auth: InvalidAuthData,
				Entity: &types.Entity{
					ID: proto.String("entity1"),
				},
			},
			wantErr:  ErrUnauthenticated,
			readonly: false,
		},
		{
			// Fails, token lacks capabilities
			req: pb.EntityRequest{
				Auth: EmptyAuthData,
				Entity: &types.Entity{
					ID: proto.String("entity1"),
				},
			},
			wantErr:  ErrRequestorUnqualified,
			readonly: false,
		},
		{
			// Fails, internal write error
			req: pb.EntityRequest{
				Auth: ValidAuthData,
				Entity: &types.Entity{
					ID: proto.String("load-error"),
				},
			},
			wantErr:  ErrInternal,
			readonly: false,
		},
		{
			// Fails, unknown user
			req: pb.EntityRequest{
				Auth: ValidAuthData,
				Entity: &types.Entity{
					ID: proto.String("does-not-exist"),
				},
			},
			wantErr:  ErrDoesNotExist,
			readonly: false,
		},
	}

	for i, c := range cases {
		s := newServer(t)
		initTree(t, s)
		s.readonly = c.readonly
		if _, err := s.EntityUnlock(context.Background(), &c.req); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}

func TestEntityGroups(t *testing.T) {
	cases := []struct {
		req     pb.EntityRequest
		wantErr error
	}{
		{
			// Works
			req: pb.EntityRequest{
				Entity: &types.Entity{
					ID: proto.String("entity1"),
				},
			},
			wantErr: nil,
		},
		{
			// Fails, entity does not exist
			req: pb.EntityRequest{
				Entity: &types.Entity{
					ID: proto.String("does-not-exist"),
				},
			},
			wantErr: ErrDoesNotExist,
		},
		{
			// Fails, entity can't be loaded
			req: pb.EntityRequest{
				Entity: &types.Entity{
					ID: proto.String("load-error"),
				},
			},
			wantErr: ErrInternal,
		},
	}

	for i, c := range cases {
		s := newServer(t)
		initTree(t, s)
		res, err := s.EntityGroups(context.Background(), &c.req)
		if err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
		t.Log(res)
	}
}
