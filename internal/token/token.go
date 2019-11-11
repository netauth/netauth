package token

import (
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"
)

// A Factory returns a token service when called.
type Factory func() (Service, error)

// The Service type defines the required interface for the Token
// Service.  The service must generate tokens, and be able to validate
// them.
type Service interface {
	Generate(Claims, Config) (string, error)
	Validate(string) (Claims, error)
}

// The Config struct contains information that should be used when
// generating a token.
type Config struct {
	Lifetime  time.Duration
	Issuer    string
	IssuedAt  time.Time
	NotBefore time.Time
}

var (
	services map[string]Factory
	logger   = hclog.L().Named("token")
)

func init() {
	services = make(map[string]Factory)
}

// New returns an initialized token service based on the value of the
// --token_impl flag.
func New() (Service, error) {
	backend := viper.GetString("token.backend")
	if backend == "" && len(services) == 1 {
		backend = GetBackendList()[0]
		logger.Warn("No implementation specified, selecting single option", "backend", backend)
	}

	t, ok := services[backend]
	if !ok {
		return nil, ErrUnknownTokenService
	}
	return t()
}

// Register is called by implementations to register ServiceFactory
// functions.
func Register(name string, impl Factory) {
	if _, ok := services[name]; ok {
		// Already registered
		return
	}
	services[name] = impl
}

// GetBackendList returns a []string of implementation names.
func GetBackendList() []string {
	var l []string

	for b := range services {
		l = append(l, b)
	}

	return l
}

// GetConfig returns a struct containing the configuration for the
// token service to use while issuing tokens.
func GetConfig() Config {
	return Config{
		Lifetime:  viper.GetDuration("token.lifetime"),
		IssuedAt:  time.Now(),
		NotBefore: time.Now(),
	}
}
