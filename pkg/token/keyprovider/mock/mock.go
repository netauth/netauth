package mock

import (
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/mock"

	"github.com/netauth/netauth/pkg/token/keyprovider"
)

// Provider is a thin wrapper on top of the mock library to allow
// mocking values.
type Provider struct {
	mock.Mock
}

// New returns a new mocked key provider.
func New(hclog.Logger) (keyprovider.KeyProvider, error) {
	return &Provider{}, nil
}

// Provider just wraps the underlying mock calls.
func (mp Provider) Provide(mech, use string) ([]byte, error) {
	args := mp.Called(mech, use)
	return args.Get(0).([]byte), args.Error(1)
}
