package bcrypt

import (
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"

	"github.com/netauth/netauth/internal/crypto"
)

func TestEncryptDecrypt(t *testing.T) {
	secret := "foo"

	viper.Set("crypto.bcrypt.cost", 0)
	e, err := New(hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}
	hash, err := e.SecureSecret(secret)
	if err != nil {
		t.Errorf("Bcrypt error: %s", err)
	}

	if err := e.VerifySecret(secret, hash); err != nil {
		t.Log(hash)
		t.Errorf("Bcrypt error: %s", err)
	}
}

func TestBadDecode(t *testing.T) {
	viper.Set("crypto.bcrypt.cost", 0)
	e, err := New(hclog.NewNullLogger())
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
	// errors out, whereas a high cost might cause the algorithm
	// to either draw down the random pool, or just lock up the
	// machine spinning.
	viper.Set("crypto.bcrypt.cost", 250)
	secret := "foo"

	e, err := New(hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}
	_, err = e.SecureSecret(secret)
	if err != crypto.ErrInternalError {
		t.Errorf("Bcrypt error: %s", err)
	}
}

// This is purely for maintaining 100% statement coverage.
func TestCB(t *testing.T) {
	cb()
}
