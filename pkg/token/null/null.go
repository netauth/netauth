// Package null implements a testing shim for testing the token system
// and some higher level components.  It is a prime candidate to be
// replaced with a mocked version of the interface, since unlike some
// other shims it is used exclusively for testing.
package null

import (
	"encoding/json"

	"github.com/hashicorp/go-hclog"

	"github.com/netauth/netauth/pkg/token"
	"github.com/netauth/netauth/pkg/token/keyprovider"
)

var (
	// ValidToken is a valid token for entity1 which has
	// GLOBAL_ROOT capability.
	ValidToken = "{\"EntityID\":\"valid\",\"Capabilities\":[\"GLOBAL_ROOT\"]}"

	// ValidEmptyToken is a valid token, but contains no
	// capabilities.
	ValidEmptyToken = "{\"EntityID\":\"valid\",\"Capabilities\":[]}"

	// InvalidToken is a token which will always return a in
	// ErrTokenInvalid error.
	InvalidToken = "invalid"
)

// Service binds the methods of the null token implementation.
type Service struct{}

// New returns a new token service
func New(_ hclog.Logger, _ keyprovider.KeyProvider) *Service {
	return &Service{}
}

// Generate generates a token with some quirks.  If the claims passed
// in requests an EntityID of "invalid-token" then an invalid token
// will be issued.  If the ID is "token-issue-error" then an
// InternalError will be returned.  For all other values a valid token
// will be returned.
func (s *Service) Generate(claims token.Claims, config token.Config) (string, error) {
	if claims.EntityID == "invalid-token" {
		return "invalid", nil
	} else if claims.EntityID == "token-issue-error" {
		return "", token.ErrInternalError
	}

	// We do this unchecked as this function will only ever see
	// prepared values, and we synthetically check for an issuance
	// error above.
	st, _ := json.Marshal(claims)
	return string(st), nil
}

// Validate checks a token that's been provided.  The token is simply
// deserialized, there is no validation performed.  Do not use this in
// production, it is provided as a testing aid only.
func (s *Service) Validate(t string) (token.Claims, error) {
	var c token.Claims
	if err := json.Unmarshal([]byte(t), &c); err != nil {
		return token.Claims{}, token.ErrTokenInvalid
	}
	return c, nil
}
