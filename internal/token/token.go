package token

import (
	"flag"
	"log"
	"time"
)

// A ServiceFactory returns a token service when called.
type ServiceFactory func() (Service, error)

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
	Renewals  int
	Issuer    string
	IssuedAt  time.Time
	NotBefore time.Time
}

var (
	services map[string]ServiceFactory

	impl     = flag.String("token_impl", "", "Token implementation to use")
	lifetime = flag.Duration("token_lifetime", time.Hour*10, "Token lifetime")
	renewals = flag.Int("token_renewals", 5, "Maximum number of times the token may be renewed")
)

func init() {
	services = make(map[string]ServiceFactory)
}

// New returns an initialized token service based on the value of the
// --token_impl flag.
func New() (Service, error) {
	if *impl == "" && len(services) == 1 {
		log.Println("Warning: No token implementation selected, using only registered option...")
		*impl = GetBackendList()[0]
	}

	t, ok := services[*impl]
	if !ok {
		return nil, ErrUnknownTokenService
	}
	return t()
}

// Register is called by implementations to register ServiceFactory
// functions.
func Register(name string, impl ServiceFactory) {
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
		Lifetime:  *lifetime,
		Renewals:  *renewals,
		IssuedAt:  time.Now(),
		NotBefore: time.Now(),
	}
}
