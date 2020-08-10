package null

import (
	"testing"

	"github.com/hashicorp/go-hclog"

	"github.com/netauth/netauth/internal/token"
)

func TestGenerate(t *testing.T) {
	tkn := New(hclog.NewNullLogger())

	cases := []struct {
		claims    token.Claims
		wantToken string
		wantErr   error
	}{
		{token.Claims{EntityID: "invalid-token"}, "invalid", nil},
		{token.Claims{EntityID: "token-issue-error"}, "", token.ErrInternalError},
		{token.Claims{EntityID: "valid"}, "{\"EntityID\":\"valid\",\"Capabilities\":null}", nil},
	}

	for i, c := range cases {
		res, err := tkn.Generate(c.claims, token.GetConfig())
		if err != c.wantErr || res != c.wantToken {
			t.Errorf("%d: Got %v, %v; Want %v, %v", i, err, res, c.wantErr, c.wantToken)
		}
	}
}

func TestValidate(t *testing.T) {
	tkn := New(hclog.NewNullLogger())

	if _, err := tkn.Validate("{\"EntityID\":\"valid\",\"Capabilities\":null}"); err != nil {
		t.Errorf("Couldn't validate real token: %v", err)
	}

	if _, err := tkn.Validate(""); err != token.ErrTokenInvalid {
		t.Error("Validated invalid token")
	}
}
