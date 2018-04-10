package token

import (
	"errors"
	"flag"
	"time"
)

type TokenServiceFactory func() (TokenService, error)

type TokenService interface {
	Generate(Claims, TokenConfig) (string, error)
	Validate(string) (Claims, error)
}

type TokenConfig struct {
	Lifetime  time.Duration
	Renewals  int
	Issuer    string
	IssuedAt  time.Time
	NotBefore time.Time
}

type Claims struct {
	EntityID     string
	Capabilities []string
	RenewalsLeft int
}

var (
	services map[string]TokenServiceFactory

	impl     = flag.String("token_impl", "", "Token implementation to use")
	lifetime = flag.Duration("token_lifetime", time.Hour*10, "Token lifetime")
	renewals = flag.Int("token_renewals", 5, "Maximum number of times the token may be renewed")

	NO_SUCH_TOKENSERVICE = errors.New("No token service with that name exists")
	KEY_UNAVAILABLE      = errors.New("Keys are not available")
	NO_GENERATE_KEYS     = errors.New("Key generation is disabled!")
)

func init() {
	services = make(map[string]TokenServiceFactory)
}

func New(impl string) (TokenService, error) {
	t, ok := services[impl]
	if !ok {
		return nil, NO_SUCH_TOKENSERVICE
	}
	return t()
}

func RegisterService(name string, impl TokenServiceFactory) {
	if _, ok := services[name]; ok {
		// Already registered
		return
	}
	services[name] = impl
}

func ListImpls() []string {
	l := []string{}

	for b, _ := range services {
		l = append(l, b)
	}

	return l
}

func GetConfig() TokenConfig {
	return TokenConfig{
		Lifetime: *lifetime,
		Renewals: *renewals,
	}
}
