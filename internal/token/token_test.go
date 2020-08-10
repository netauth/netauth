package token

import (
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
)

type dummyTokenService struct{}

func (*dummyTokenService) Generate(Claims, Config) (string, error) { return "", nil }
func (*dummyTokenService) Validate(string) (Claims, error)         { return Claims{}, nil }
func newDummyTokenService(_ hclog.Logger) (Service, error)         { return new(dummyTokenService), nil }

func TestRegister(t *testing.T) {
	services = make(map[string]Factory)

	Register("dummy", newDummyTokenService)
	if len(services) != 1 {
		t.Error("Service factory failed to register")
	}

	Register("dummy", newDummyTokenService)
	if len(services) != 1 {
		t.Error("A duplicate TokenService was registered")
	}
}

func TestNewKnown(t *testing.T) {
	services = make(map[string]Factory)

	Register("dummy", newDummyTokenService)

	x, err := New("dummy")
	if err != nil {
		t.Error(err)
	}

	if _, ok := x.(*dummyTokenService); !ok {
		t.Error("Returned implementation is incorrect")
	}
}

func TestNewUnknown(t *testing.T) {
	services = make(map[string]Factory)

	if x, err := New("unknown"); x != nil || err != ErrUnknownTokenService {
		t.Error("Undefined error behavior")
	}
}

func TestGetConfig(t *testing.T) {
	c := GetConfig()
	if c.Lifetime != time.Minute*5 {
		t.Error("Config contains incorrect values")
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

func TestSetLifetime(t *testing.T) {
	SetLifetime(time.Second * 42)
	if lifetime != time.Second*42 {
		t.Error("Wrong duration")
	}
}
