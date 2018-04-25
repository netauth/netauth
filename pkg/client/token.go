package client

import (
	"github.com/NetAuth/NetAuth/internal/token"
)

func (n *netAuthClient) InspectToken(t string) (token.Claims, error) {
	return n.tokenService.Validate(t)
}
