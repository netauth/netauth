package keyprovider

import (
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

type testProvider struct{}

func newTestProvider(hclog.Logger) (KeyProvider, error)        { return testProvider{}, nil }
func (tp testProvider) Provide(string, string) ([]byte, error) { return []byte{}, nil }

func TestRegister(t *testing.T) {
	providers = make(map[string]Factory)

	Register("test", newTestProvider)
	if len(providers) != 1 {
		t.Error("KeyProvider factory failed to register")
	}

	Register("test", newTestProvider)
	if len(providers) != 1 {
		t.Error()
	}
}

func TestNewKnown(t *testing.T) {
	providers = make(map[string]Factory)

	Register("test", newTestProvider)

	_, err := New("test")
	assert.Nil(t, err)
}

func TestNewUnknown(t *testing.T) {
	providers = make(map[string]Factory)

	_, err := New("unknown")
	assert.Equal(t, err, ErrUnknownKeyProvider)
}

func TestSetParentLogger(t *testing.T) {
	lb = nil
	assert.Nil(t, lb)
	l := hclog.NewNullLogger()
	SetParentLogger(l)
	assert.NotNil(t, log())
}
