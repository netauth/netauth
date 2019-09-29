package token

import (
	"testing"

	pb "github.com/NetAuth/Protocol"
)

func TestHasCapability(t *testing.T) {
	cases := []struct {
		c     []pb.Capability
		check pb.Capability
		want  bool
	}{
		{[]pb.Capability{pb.Capability_CREATE_ENTITY}, pb.Capability_CREATE_ENTITY, true},
		{[]pb.Capability{pb.Capability_GLOBAL_ROOT}, pb.Capability_CREATE_ENTITY, true},
		{[]pb.Capability{pb.Capability_DESTROY_ENTITY}, pb.Capability_CREATE_ENTITY, false},
	}

	for i, c := range cases {
		claims := Claims{
			Capabilities: c.c,
		}

		if claims.HasCapability(c.check) != c.want {
			t.Errorf("%d: Got %t Want %t", i, claims.HasCapability(c.check), c.want)
		}
	}
}
