package crypto

import (
	"testing"

	"github.com/hashicorp/go-hclog"
)

type dummyCrypto struct{}

func (*dummyCrypto) SecureSecret(_ string) (string, error) { return "", nil }
func (*dummyCrypto) VerifySecret(_, _ string) error        { return nil }
func dummyCryptoFactory(_ hclog.Logger) (EMCrypto, error)  { return new(dummyCrypto), nil }
func dummyCryptoCallback()                                 { Register("dummy", dummyCryptoFactory) }

func TestRegister(t *testing.T) {
	backends = make(map[string]Factory)

	Register("dummy", dummyCryptoFactory)
	if len(backends) != 1 {
		t.Error("Engine wasn't registered")
	}

	Register("dummy", dummyCryptoFactory)
	if len(backends) != 1 {
		t.Error("Wrong number of engines")
	}
}

func TestNewKnown(t *testing.T) {
	backends = make(map[string]Factory)

	Register("dummy", dummyCryptoFactory)

	x, err := New("dummy")
	if err != nil {
		t.Error(err)
	}

	if _, ok := x.(*dummyCrypto); !ok {
		t.Error("Something that isn't a crypto engine came out...")
	}
}

func TestNewUnknown(t *testing.T) {
	backends = make(map[string]Factory)
	x, err := New("foobar")
	if x != nil && err != ErrUnknownCrypto {
		t.Error(err)
	}
}

func TestSetParentLogger(t *testing.T) {
	lb = nil

	l := hclog.NewNullLogger()
	SetParentLogger(l)
	if log() == nil {
		t.Error("log was not set")
	}
}

func TestLogParentUnset(t *testing.T) {
	lb = nil

	if log() == nil {
		t.Error("auto log was not aquired")
	}
}

func TestRegisterCallback(t *testing.T) {
	callbacks = nil
	RegisterCallback(dummyCryptoCallback)
	if len(callbacks) != 1 {
		t.Error("Callback not registered")
	}
}

func TestDoCallbacks(t *testing.T) {
	callbacks = nil
	called := false

	testCB := func() {
		called = true
	}

	RegisterCallback(testCB)
	DoCallbacks()

	if !called {
		t.Error("Callback was not called")
	}
}
