package keyprovider

import (
	"github.com/hashicorp/go-hclog"
)

// Factory returns a KeyProvider to the caller.
type Factory func(hclog.Logger) (KeyProvider, error)

// KeyProvider provides an interface for obtaining keys from various
// sources.
type KeyProvider interface {
	Provide(string, string) ([]byte, error)
}

var (
	lb        hclog.Logger
	providers map[string]Factory
)

func init() {
	providers = make(map[string]Factory)
}

// Register is called by implementations during early init to register
// themselves.
func Register(name string, impl Factory) {
	if _, ok := providers[name]; ok {
		return
	}
	providers[name] = impl
}

// New returns a KeyProvider initialized with the specified logger.
func New(b string) (KeyProvider, error) {
	p, ok := providers[b]
	if !ok {
		return nil, ErrUnknownKeyProvider
	}
	return p(log())
}

// SetParentLogger sets the parent logger for this instance.
func SetParentLogger(l hclog.Logger) {
	lb = l.Named("keyprovider")
}

// log is a convenience function that will return a null logger if a
// parent logger has not been specified, mostly useful for tests.
func log() hclog.Logger {
	if lb == nil {
		lb = hclog.NewNullLogger()
	}
	return lb
}
