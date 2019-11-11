package netauth

var (
	tokenCacheFactories map[string]TokenCacheFactory
)

func init() {
	tokenCacheFactories = make(map[string]TokenCacheFactory)
}

// The TokenCache is an interface that allows many possible token
// caches to be used to improve token efficiency.
type TokenCache interface {
	PutToken(string, string) error
	GetToken(string) (string, error)
	DelToken(string) error
}

// A TokenCacheFactory is a function that returns a TokenCache.  This
// should generally not be an exported function in the implementation,
// and is instead registered here on import.
type TokenCacheFactory func() (TokenCache, error)

// RegisterTokenCacheFactory can be used by token caches on import to
// register themselves.  Names must be unique.
func RegisterTokenCacheFactory(name string, f TokenCacheFactory) {
	if _, ok := tokenCacheFactories[name]; ok {
		// A cache already exists, cowardly refusing to
		// overwrite it...
		return
	}
	tokenCacheFactories[name] = f
}

// NewTokenCache returns a token cache that's setup and ready for use.
// While its not expected that this function will be called outside of
// this package, there are some rare circumstances that it would be
// useful to initialize a cache for private use.
func NewTokenCache(name string) (TokenCache, error) {
	f, ok := tokenCacheFactories[name]
	if !ok {
		return nil, ErrUnknownCache
	}
	return f()
}
