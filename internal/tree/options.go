package tree

import (
	"github.com/hashicorp/go-hclog"

	"github.com/netauth/netauth/internal/crypto"
)

func WithStorage(d DB) Option {
	return func(m *Manager) { m.db = d }
}

func WithCrypto(c crypto.EMCrypto) Option {
	return func(m *Manager) { m.crypto = c }
}

func WithLogger(l hclog.Logger) Option {
	return func(m *Manager) { m.log = l.Named("tree") }
}
