package nocrypto

import (
	"testing"

	"github.com/NetAuth/NetAuth/internal/crypto"
)

func TestSecureSecret(t *testing.T) {
	e := New()
	s := "foo"
	h, err := e.SecureSecret("foo")
	if h != s && err != nil {
		t.Errorf("NoCrypto wtf!? %s != %s | %s", h, s, err)
	}
}

func TestSecureSecretBadAuth(t *testing.T) {
	e := New()
	s := "foo"
	h := "bar"

	if err := e.VerifySecret(s, h); err != crypto.ErrAuthorizationFailure {
		t.Error(err)
	}
}

func TestVerifySecret(t *testing.T) {
	e := New()
	h := "foo"
	s := "foo"
	if err := e.VerifySecret(s, h); err != nil {
		t.Errorf("NoCrypto wtf!? %s", err)
	}
}
