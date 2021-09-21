package interface_test

import (
	"context"
	"testing"

	"github.com/netauth/netauth/internal/crypto"
)

func TestValidateSecret(t *testing.T) {
	ctxt := context.Background()
	m, ctx := newTreeManager(t)

	addEntity(t, ctx)

	if err := m.ValidateSecret(ctxt, "entity1", "entity1"); err != nil {
		t.Error(err)
	}

	if err := m.ValidateSecret(ctxt, "entity1", "password"); err != crypto.ErrAuthorizationFailure {
		t.Error(err)
	}
}
