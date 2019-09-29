package nocrypto

import (
	"testing"

	"github.com/NetAuth/NetAuth/internal/crypto"
)

func TestSecureSecret(t *testing.T) {
	e, err := New()
	if err != nil {
		t.Fatal(err)
	}
	s := "foo"
	h, err := e.SecureSecret("foo")
	if h != s && err != nil {
		t.Errorf("NoCrypto wtf!? %s != %s | %s", h, s, err)
	}

	h, err = e.SecureSecret("return-error")
	if h != "" || err != crypto.ErrInternalError {
		t.Errorf("Trigger secret failed to trigger an error")
	}
}

func TestSecureSecretBadAuth(t *testing.T) {
	e, err := New()
	if err != nil {
		t.Fatal(err)
	}
	s := "foo"
	h := "bar"

	if err := e.VerifySecret(s, h); err != crypto.ErrAuthorizationFailure {
		t.Error(err)
	}
}

func TestVerifySecret(t *testing.T) {
	e, err := New()
	if err != nil {
		t.Fatal(err)
	}
	h := "foo"
	s := "foo"
	if err := e.VerifySecret(s, h); err != nil {
		t.Errorf("NoCrypto wtf!? %s", err)
	}
}
