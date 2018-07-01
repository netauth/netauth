package health

import "testing"

func TestHealth(t *testing.T) {
	s := []struct {
		set  bool
		want bool
	}{
		{false, false},
		{true, true},
	}

	for _, c := range s {
		if c.set {
			SetGood()
		} else {
			SetBad()
		}

		if Get() != c.want {
			t.Error("Bad health status")
		}
	}
}
