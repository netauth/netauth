package client

import (
	"context"

	"github.com/netauth/netauth/internal/token"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/NetAuth/Protocol"
)

// GetToken is identical to Authenticate except on success it will
// return a token which can be used to authorize additional later
// requests.
func (n *NetAuthClient) GetToken(entity, secret string) (string, error) {
	// See if we have a local copy first.
	t, err := n.getTokenFromStore(entity)
	if err == nil {
		var err error
		if _, err = n.InspectToken(t); err == nil {
			n.log.Debug("Using cached token")
			return t, nil
		}
		n.log.Debug("Not using cached token", "error", err)
	} else {
		n.log.Debug("Could not retrieve cached token", "error", err)
	}

	if secret == "" {
		return "", ErrTokenUnavailable
	}

	request := pb.NetAuthRequest{
		Entity: &pb.Entity{
			ID:     &entity,
			Secret: &secret,
		},
		Info: clientInfo(),
	}
	tokenResult, err := n.c.GetToken(context.Background(), &request)
	if status.Code(err) != codes.OK {
		return "", err
	}

	t = tokenResult.GetToken()
	if err := n.tokenStore.DestroyToken(entity); err != nil {
		return "", err
	}
	err = n.putTokenInStore(entity, t)
	return t, err
}

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
