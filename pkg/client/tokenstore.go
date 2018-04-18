package client

import (
	"errors"
	"flag"
)

// This file implements a plugin token storage system, it is
// overbuilt.  This is really not necessary but it makes it easier to
// build the client library like all the other plugin systems in
// NetAuth.  In reality the only sane backends are likely memory and
// file.

type TokenStore interface {
	StoreToken(string, string) error
	GetToken(string) (string, error)
}

var (
	backends map[string]TokenStore
	impl     = flag.String("tokenstore", "", "Token storage system")

	NoSuchTokenStore = errors.New("No token store with that name exists!")
	TokenUnavailable = errors.New("The stored token is unavailable")
)

func init() {
	backends = make(map[string]TokenStore)
	Register("memory", &memTokenStore{})
}

// Mechanism to register token storage systems
func Register(name string, impl TokenStore) {
	if _, ok := backends[name]; ok {
		// Already registered
		return
	}
	backends[name] = impl
}

func getTokenStore() (TokenStore, error) {
	// If nothing was specified select the only backend
	// registered, if more were registered than the user has to
	// make an stated choice.
	if *impl == "" && len(backends) == 1 {
		for b, _ := range backends {
			*impl = b
			break
		}
	}

	if t, ok := backends[*impl]; ok {
		return t, nil
	}
	return nil, NoSuchTokenStore
}

// Exposed functions to store and retrieve the tokens
func (n *netAuthClient) storeToken(name, token string) error {
	return n.tokenStore.StoreToken(name, token)
}

func (n *netAuthClient) getTokenFromStore(name string) (string, error) {
	t, err := n.tokenStore.GetToken(name)
	if err != nil {
		return "", err
	}
	if t == "" {
		return "", TokenUnavailable
	}
	return n.tokenStore.GetToken(name)
}

// Basic in memory token store
type memTokenStore struct {
	token string
}

func (m *memTokenStore) StoreToken(name, token string) error {
	m.token = token
	return nil
}

func (m *memTokenStore) GetToken(name string) (string, error) {
	return m.token, nil
}
