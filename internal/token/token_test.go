package token

import (
	"testing"

	"github.com/spf13/viper"
)

type dummyTokenService struct{}

func (*dummyTokenService) Generate(Claims, Config) (string, error) { return "", nil }
func (*dummyTokenService) Validate(string) (Claims, error)         { return Claims{}, nil }
func newDummyTokenService() (Service, error)                       { return new(dummyTokenService), nil }

func TestRegister(t *testing.T) {
	services = make(map[string]Factory)

	Register("dummy", newDummyTokenService)
	if l := GetBackendList(); len(l) != 1 || l[0] != "dummy" {
		t.Error("Service factory failed to register")
	}

	Register("dummy", newDummyTokenService)
	if l := GetBackendList(); len(l) != 1 {
		t.Error("A duplicate TokenService was registered")
	}
}

func TestNewKnown(t *testing.T) {
	services = make(map[string]Factory)

	Register("dummy", newDummyTokenService)

	viper.Set("token.backend", "dummy")
	x, err := New()
	if err != nil {
		t.Error(err)
	}

	if _, ok := x.(*dummyTokenService); !ok {
		t.Error("Returned implementation is incorrect")
	}
}

func TestNewUnspecified(t *testing.T) {
	services = make(map[string]Factory)

	Register("dummy", newDummyTokenService)

	viper.Set("token.backend", "")
	x, err := New()
	if err != nil {
		t.Error(err)
	}

	if _, ok := x.(*dummyTokenService); !ok {
		t.Error("Returned implementation is incorrect")
	}

}

func TestNewUnknown(t *testing.T) {
	services = make(map[string]Factory)

	viper.Set("token.backend", "unknown")
	if x, err := New(); x != nil || err != ErrUnknownTokenService {
		t.Error("Undefined error behavior")
	}
}

func TestGetConfig(t *testing.T) {
	c := GetConfig()
	if c.Lifetime != viper.GetDuration("token.lifetime") || c.Renewals != viper.GetInt("token.renewals") {
		t.Error("Config contains incorrect values")
	}
}
