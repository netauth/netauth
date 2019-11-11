package netauth

import (
	"testing"
)

type dummyCache struct{}

func newDummyCache() (TokenCache, error)            { return &dummyCache{}, nil }
func (*dummyCache) PutToken(string, string) error   { return nil }
func (*dummyCache) GetToken(string) (string, error) { return "", nil }
func (*dummyCache) DelToken(string) error           { return nil }

func TestRegisterTokenCacheFactory(t *testing.T) {
	// Just call it twice to verify that it works.  This can only
	// explode if the map isn't initialized correctly and that
	// will generate a panic that will subsequently fail this test
	// case.
	RegisterTokenCacheFactory("dummy", newDummyCache)
	RegisterTokenCacheFactory("dummy", newDummyCache)
}

func TestNewTokenCache(t *testing.T) {
	RegisterTokenCacheFactory("dummy", newDummyCache)

	if _, err := NewTokenCache("dummy"); err != nil {
		t.Errorf("Unexpected error while initializing cache: %s", err)
	}

	if _, err := NewTokenCache("does-not-exist"); err != ErrUnknownCache {
		t.Errorf("Unexpected error while initializing non-existent cache: %s", err)
	}
}
