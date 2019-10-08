package rpc2

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"

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
		{UnprivilegedContext, ErrRequestorUnqualified},
		{UnauthenticatedContext, ErrMalformedRequest},
		{InvalidAuthContext, ErrUnauthenticated},
	}

	for i, c := range cases {
		s := newServer(t)

		if err := s.checkToken(c.ctx, types.Capability_GLOBAL_ROOT); err != c.wantErr {
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
