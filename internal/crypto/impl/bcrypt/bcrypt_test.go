package bcrypt

import (
	"testing"

	"github.com/NetAuth/NetAuth/internal/crypto"
)

func TestEncryptDecrypt(t *testing.T) {
	secret := "foo"

	*cost = 0
	e := New()
	hash, err := e.SecureSecret(secret)
	t.Log("Testing GenerateFromPassword")
	if err != nil {
		t.Errorf("Bcrypt error: %s", err)
	}

	t.Log("Testing CompareHashAndPassword")
	if err := e.VerifySecret(secret, hash); err != nil {
		t.Log(hash)
		t.Errorf("Bcrypt error: %s", err)
	}
}

func TestBadDecode(t *testing.T) {
	*cost = 0
	e := New()

	if err := e.VerifySecret("", ""); err != crypto.ErrAuthorizationFailure {
		t.Errorf("Bad crypto error: %s", err)
	}
}
