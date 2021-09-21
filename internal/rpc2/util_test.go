package rpc2

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"

	"github.com/netauth/netauth/internal/token"

	types "github.com/netauth/protocol"
)

func TestGetCapabilitiesForEntity(t *testing.T) {
	s := newServer(t)
	initTree(t, s.Manager)

	s.CreateGroup(context.Background(), "lockout", "", "", -1)
	s.AddEntityToGroup(context.Background(), "admin", "lockout")
	s.SetGroupCapability2(context.Background(), "lockout", types.Capability_LOCK_ENTITY.Enum())

	caps := s.getCapabilitiesForEntity(context.Background(), "admin")
	if len(caps) != 2 {
		t.Error("Not all caps were found", caps)
	}
}

func TestCheckToken(t *testing.T) {
	cases := []struct {
		ctx     context.Context
		wantErr error
	}{
		{PrivilegedContext, nil},
		{UnauthenticatedContext, ErrMalformedRequest},
		{InvalidAuthContext, ErrUnauthenticated},
	}

	for i, c := range cases {
		s := newServer(t)

		if _, err := s.checkToken(c.ctx); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}

func TestGetSingleStringFromMetaData(t *testing.T) {
	cases := []struct {
		ctx     context.Context
		wantRes string
	}{
		{metadata.NewIncomingContext(context.Background(), metadata.Pairs("key", "value")), "value"},
		{metadata.NewIncomingContext(context.Background(), metadata.Pairs("key", "value", "key", "value2")), ""},
		{context.Background(), ""},
	}

	for i, c := range cases {
		if res := getSingleStringFromMetadata(c.ctx, "key"); res != c.wantRes {
			t.Errorf("%d: Got %v; Want %v", i, res, c.wantRes)
		}
	}
}

func TestGetClientName(t *testing.T) {
	cases := []struct {
		ctx     context.Context
		wantRes string
	}{
		{metadata.NewIncomingContext(context.Background(), metadata.Pairs("client-name", "test")), "test"},
		{context.Background(), "BOGUS_CLIENT"},
	}

	for i, c := range cases {
		if r := getClientName(c.ctx); r != c.wantRes {
			t.Errorf("%d: Got %v; Want %v", i, r, c.wantRes)
		}
	}
}

func TestGetServiceName(t *testing.T) {
	cases := []struct {
		ctx     context.Context
		wantRes string
	}{
		{metadata.NewIncomingContext(context.Background(), metadata.Pairs("service-name", "test")), "test"},
		{context.Background(), "BOGUS_SERVICE"},
	}

	for i, c := range cases {
		if r := getServiceName(c.ctx); r != c.wantRes {
			t.Errorf("%d: Got %v; Want %v", i, r, c.wantRes)
		}
	}
}

func TestGetTokenClaims(t *testing.T) {
	cases := []struct {
		ctx       context.Context
		wantEmpty bool
	}{
		{context.WithValue(context.Background(), claimsContextKey{}, token.Claims{EntityID: "foo"}), false},
		{context.Background(), true},
	}

	for i, c := range cases {
		res := getTokenClaims(c.ctx)
		if res.EntityID == "" && !c.wantEmpty {
			t.Errorf("%d: Got Empty claims when shouldn't have", i)
		}
	}
}

func TestManageByMembership(t *testing.T) {
	cases := []struct {
		id      string
		g       types.Group
		wantRes bool
	}{
		{
			id: "entity1",
			g: types.Group{
				Name: proto.String("group2"),
			},
			wantRes: true,
		},
		{
			id: "entity1",
			g: types.Group{
				Name: proto.String("group1"),
			},
			wantRes: false,
		},
		{
			id: "entity1",
			g: types.Group{
				Name: proto.String("load-error"),
			},
			wantRes: false,
		},
		{
			id: "load-error",
			g: types.Group{
				Name: proto.String("group2"),
			},
			wantRes: false,
		},
		{
			id: "unprivileged",
			g: types.Group{
				Name: proto.String("group2"),
			},
			wantRes: false,
		},
	}

	for i, c := range cases {
		s := newServer(t)
		initTree(t, s.Manager)

		if got := s.manageByMembership(context.Background(), c.id, &c.g); got != c.wantRes {
			t.Errorf("%d: Got %v; Want %v", i, got, c.wantRes)
		}
	}
}

func TestMutablePrequisitesMet(t *testing.T) {
	cases := []struct {
		ro      bool
		ctx     context.Context
		cap     types.Capability
		wantErr error
	}{
		{
			ro:      true,
			ctx:     PrivilegedContext,
			cap:     types.Capability_CREATE_ENTITY,
			wantErr: ErrReadOnly,
		},
		{
			ro:      false,
			ctx:     InvalidAuthContext,
			cap:     types.Capability_CREATE_ENTITY,
			wantErr: ErrUnauthenticated,
		},
		{
			ro:      false,
			ctx:     UnprivilegedContext,
			cap:     types.Capability_CREATE_ENTITY,
			wantErr: ErrRequestorUnqualified,
		},
		{
			ro:      false,
			ctx:     PrivilegedContext,
			cap:     types.Capability_CREATE_ENTITY,
			wantErr: nil,
		},
	}

	for i, c := range cases {
		s := newServer(t)
		initTree(t, s.Manager)
		s.readonly = c.ro

		assert.Equalf(t, c.wantErr, s.mutablePrequisitesMet(c.ctx, c.cap), "Test Number %d", i)
	}
}
