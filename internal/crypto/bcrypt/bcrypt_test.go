package bcrypt

import (
	"testing"

	"github.com/NetAuth/NetAuth/internal/crypto"
)

func TestEncryptDecrypt(t *testing.T) {
	secret := "foo"

	*cost = 0
	e, err := New()
	if err != nil {
		t.Fatal(err)
	}
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
	e, err := New()
	if err != nil {
		t.Fatal(err)
	}

	if err := e.VerifySecret("", ""); err != crypto.ErrAuthorizationFailure {
		t.Errorf("Bad crypto error: %s", err)
	}
}

func TestCostTooHigh(t *testing.T) {
	// This needs to be an arbitrarily high number above maxCost.
	// The bcrypt library doesn't clamp for high cost and instead
	// errors out, wheras a high cost might cause the algorithm to
	// either draw down the random pool, or just lock up the
	// machine spinning.
	*cost = 250
	secret := "foo"

	e, err := New()
	if err != nil {
		t.Fatal(err)
	}
	_, err = e.SecureSecret(secret)
	t.Log("Testing GenerateFromPassword")
	if err != crypto.ErrInternalError {
		t.Errorf("Bcrypt error: %s", err)
	}
}
