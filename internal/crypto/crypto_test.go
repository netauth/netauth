package crypto

import "testing"

type dummyCrypto struct {}
func (*dummyCrypto) SecureSecret(_ string) (string, error) { return "", nil }
func (*dummyCrypto) VerifySecret(_, _ string) error { return nil }
func dummyCryptoFactory() (EMCrypto, error) { return new(dummyCrypto), nil }

func TestRegister(t *testing.T) {
	backends = make(map[string]Factory)

	Register("dummy", dummyCryptoFactory)
	if l := GetBackendList(); len(l) != 1 && l[0] != "dummy" {
		t.Error("Engine wasn't registered")
	}

	Register("dummy", dummyCryptoFactory)
	if l := GetBackendList(); len(l) != 1 {
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
	x, err := New("unknown")
	if x != nil && err != ErrUnknownCrypto {
		t.Error(err)
	}
}
