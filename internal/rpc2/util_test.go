package rpc2

import (
	"context"
	"testing"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/metadata"

	"github.com/netauth/netauth/internal/token"

	types "github.com/NetAuth/Protocol"
)

func TestGetCapabilitiesForEntity(t *testing.T) {
	s := newServer(t)
	initTree(t, s)

	s.CreateGroup("lockout", "", "", -1)
	s.AddEntityToGroup("admin", "lockout")
	s.SetGroupCapability2("lockout", types.Capability_LOCK_ENTITY.Enum())

	caps := s.getCapabilitiesForEntity("admin")
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
		initTree(t, s)

		if got := s.manageByMembership(c.id, &c.g); got != c.wantRes {
			t.Errorf("%d: Got %v; Want %v", i, got, c.wantRes)
		}
	}
}
