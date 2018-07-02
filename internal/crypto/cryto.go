package crypto

// The EMCrypto interface defines the functions that are needed to
// make a secret secure for storage and later verify a secret against
// the secured copy.
type EMCrypto interface {
	SecureSecret(string) (string, error)
	VerifySecret(string, string) error
}

// The Factory type is to be implemented by crypto implementations and
// shall be fed to the Register function.
type Factory func() EMCrypto

var (
	backends = make(map[string]Factory)
)

// New returns an initialized Crypto instance which can create and
// verify secure versions of secrets.
func New(name string) (EMCrypto, error) {
	b, ok := backends[name]
	if !ok {
		return nil, ErrUnknownCrypto
	}
	return b(), nil
}

// Register takes in a name for the engine and a function
// signature to bind to that name.
func Register(name string, newFunc Factory) {
	if _, ok := backends[name]; ok {
		// Return if the backend was already registered.
		return
	}
	backends[name] = newFunc
}

// GetBackendList returns a string list of hte backends that are
// available.
func GetBackendList() []string {
	var l []string

	for b := range backends {
		l = append(l, b)
	}

	return l
}
