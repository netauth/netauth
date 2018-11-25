package interface_test

import (
	"testing"

	"github.com/NetAuth/NetAuth/internal/crypto"
)

func TestValidateSecret(t *testing.T) {
	m, ctx := newTreeManager(t)

	addEntity(t, ctx)

	if err := m.ValidateSecret("entity1", "entity1"); err != nil {
		t.Error(err)
	}

	if err := m.ValidateSecret("entity1", "password"); err != crypto.ErrAuthorizationFailure {
		t.Error(err)
	}
}
