package client

import (
	"github.com/NetAuth/NetAuth/internal/token"
)

// InspectToken proxies through to the tokenService since the inner
// function may oneday be significantly more complicated, but hte
// function in the client should not change.
func (n *NetAuthClient) InspectToken(t string) (token.Claims, error) {
	if n.tokenService == nil {
		return token.Claims{}, token.ErrKeyUnavailable
	}
	return n.tokenService.Validate(t)
}

// DestroyToken proxies inwards to the tokenStore to shield the client
// API for the future.
func (n *NetAuthClient) DestroyToken(name string) error {
	return n.tokenStore.DestroyToken(name)
}
