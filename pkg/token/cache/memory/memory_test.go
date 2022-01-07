package memory

import (
	"testing"

	"github.com/netauth/netauth/pkg/token/cache"
)

func TestInMemoryCache(t *testing.T) {
	x, err := new()
	if err != nil {
		t.Fatal(err)
	}

	if err := x.PutToken("foo", "bar"); err != nil {
		t.Error(err)
	}

	tk, err := x.GetToken("foo")
	if err != nil || tk != "bar" {
		t.Error("Error loading token")
	}

	if err := x.DelToken("foo"); err != nil {
		t.Error(err)
	}

	tk, err = x.GetToken("foo")
	if err != cache.ErrNoCachedToken || tk != "" {
		t.Error("Incorrect response for non-existent token")
	}
}
