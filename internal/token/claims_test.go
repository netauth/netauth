package token

import "testing"

func TestHasCapability(t *testing.T) {
	cases := []struct {
		c     []string
		check string
		want  bool
	}{
		{[]string{"CREATE_ENTITY"}, "CREATE_ENTITY", true},
		{[]string{"GLOBAL_ROOT"}, "CREATE_ENTITY", true},
		{[]string{"DESTORY_ENTITY"}, "CREATE_ENTITY", false},
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
