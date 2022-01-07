package fs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/netauth/netauth/pkg/token/cache"
)

func TestPutGetDeleteOK(t *testing.T) {
	x, _ := new()

	if err := x.PutToken("foo", "bar"); err != nil {
		t.Errorf("Error putting token: %v", err)
	}

	r, err := x.GetToken("foo")
	if err != nil {
		t.Errorf("Error retrieving token: %v", err)
	}
	if r != "bar" {
		t.Errorf("Incorrect value retrieved: %s", r)
	}

	if err := x.DelToken("foo"); err != nil {
		t.Errorf("Error removing token: %v", err)
	}
}

func TestGetMissingToken(t *testing.T) {
	x, _ := new()

	if _, err := x.GetToken("does-not-exist"); err != cache.ErrNoCachedToken {
		t.Errorf("Incorrect error for bogus owner: %v", err)
	}
}

func TestGetTokenBadPath(t *testing.T) {
	x, _ := new()

	rx, ok := x.(*fsCache)
	if !ok {
		t.Fatal("Non fscache in fscache implementation!?")
	}

	if err := os.MkdirAll(rx.filepathFromOwner("bad-path"), 0755); err != nil {
		t.Fatalf("Test setup failed: %v", err)
	}
	defer os.Remove(rx.filepathFromOwner("bad-path"))

	if _, err := x.GetToken("bad-path"); err == nil || err == cache.ErrNoCachedToken {
		t.Errorf("Wrong error when encountering read fail: %v", err)
	}
}

func TestDelTokenBadPath(t *testing.T) {
	x, _ := new()

	rx, ok := x.(*fsCache)
	if !ok {
		t.Fatal("Non fscache in fscache implementation!?")
	}

	p := filepath.Join(rx.filepathFromOwner("bad-path"), "remove-this")
	if err := os.MkdirAll(p, 0755); err != nil {
		t.Fatalf("Test setup failed: %v", err)
	}
	defer os.RemoveAll(p)

	if err := x.DelToken("bad-path"); err == nil {
		t.Errorf("Wrong error when encountering delete fail: %v", err)
	}
}
