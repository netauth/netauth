// Package memory implements an in-memory token cache.  This cache is
// inappropriate for use in CLI tools or call/response applications
// due to its lack of out-of-process persistence.  The goal of this
// cache is to accelerate applications such as web interfaces that
// need to persist tokens and don't want to actually send a token to
// the client.
package memory

import (
	"sync"

	"github.com/netauth/netauth/pkg/netauth"
)

type inMemoryCache struct {
	sync.RWMutex
	c map[string]string
}

func init() {
	netauth.RegisterTokenCacheFactory("memory", new)
}

func new() (netauth.TokenCache, error) {
	c := &inMemoryCache{
		c: make(map[string]string),
	}
	return c, nil
}

// PutToken stores a token in the cache.  This function may block if
// another goroutine is currently modifying the cache.
func (imc *inMemoryCache) PutToken(owner, token string) error {
	imc.Lock()
	imc.c[owner] = token
	imc.Unlock()
	return nil
}

// GetToken retrieves a token from the cache.  No promises are made
// that the returned token will be valid, unexpired, or use-able for a
// particular purpose.  No refunds if this function returns a token
// that eats your dog.
func (imc *inMemoryCache) GetToken(owner string) (string, error) {
	imc.RLock()
	res, ok := imc.c[owner]
	imc.RUnlock()
	if !ok {
		return "", netauth.ErrNoCachedToken
	}
	return res, nil
}

// DelToken should be called whenever a token is being invalidated to
// immediately remove it from the cache.  This function may block if
// another goroutine is currently modifying the cache.
func (imc *inMemoryCache) DelToken(owner string) error {
	imc.Lock()
	delete(imc.c, owner)
	imc.Unlock()
	return nil
}
