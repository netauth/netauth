package client

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// This file implements a plugin token storage system, it is
// overbuilt.  This is really not necessary but it makes it easier to
// build the client library like all the other plugin systems in
// NetAuth.  In reality the only sane backends are likely memory and
// file.

// The TokenStore is a convenient way to securely store tokens for
// entities.  Care should be taken with all implementations to avoid
// loosing security of the token, since a token attack can be
// escalated to persistent root in the right circumstances.
type TokenStore interface {
	StoreToken(string, string) error
	GetToken(string) (string, error)
	DestroyToken(string) error
}

var (
	backends map[string]TokenStore
	impl     = flag.String("tokenstore", "disk", "Token storage system")

	// ErrNoSuchTokenStore is returned in the case when the token
	// store requested does not actually exist.
	ErrNoSuchTokenStore = errors.New("no token store with that name exists")

	// ErrTokenUnavailable is returned when there is no token
	// available to be returned.
	ErrTokenUnavailable = errors.New("the stored token is unavailable")
)

func init() {
	backends = make(map[string]TokenStore)
	Register("memory", &memTokenStore{})
	Register("disk", &fsTokenStore{})
}

// Register is called by implementations to register into the token
// system.
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
		for b := range backends {
			*impl = b
			break
		}
	}

	if t, ok := backends[*impl]; ok {
		return t, nil
	}
	return nil, ErrNoSuchTokenStore
}

// Exposed functions to store and retrieve the tokens
func (n *NetAuthClient) storeToken(name, token string) error {
	return n.tokenStore.StoreToken(name, token)
}

func (n *NetAuthClient) getTokenFromStore(name string) (string, error) {
	t, err := n.tokenStore.GetToken(name)
	if err != nil {
		return "", err
	}
	if t == "" {
		return "", ErrTokenUnavailable
	}
	return n.tokenStore.GetToken(name)
}

func (n *NetAuthClient) putTokenInStore(name, token string) error {
	return n.tokenStore.StoreToken(name, token)
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

func (m *memTokenStore) DestroyToken(name string) error {
	m.token = ""
	return nil
}

// Basic filesystem token store
type fsTokenStore struct{}

func (*fsTokenStore) StoreToken(name, token string) error {
	tokenFile := filepath.Join(os.TempDir(), fmt.Sprintf("%s.%s", name, "token"))

	return ioutil.WriteFile(tokenFile, []byte(token), 0400)
}

func (*fsTokenStore) GetToken(name string) (string, error) {
	tokenFile := filepath.Join(os.TempDir(), fmt.Sprintf("%s.%s", name, "token"))

	d, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		return "", err
	}

	return string(d), nil
}

func (*fsTokenStore) DestroyToken(name string) error {
	tokenFile := filepath.Join(os.TempDir(), fmt.Sprintf("%s.%s", name, "token"))

	err := os.Remove(tokenFile)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
