package client

import (
	"github.com/NetAuth/NetAuth/internal/token"
)

func (n *NetAuthClient) InspectToken(t string) (token.Claims, error) {
	return n.tokenService.Validate(t)
}
