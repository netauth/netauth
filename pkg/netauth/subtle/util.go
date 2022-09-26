package subtle

import (
	"crypto/subtle"
	"errors"

	"golang.org/x/crypto/ssh"
)

var (
	// ErrNoMatchingKeys is returned when no matching keys are found
	ErrNoMatchingKeys = errors.New("no matching keys found")
)

// CompareSSHKeys checks if the given SSH public key is in the list of SSH public keys
func CompareSSHKeys(keys []string, test string) error {
	pubkey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(test))
	if err != nil {
		return ErrNoMatchingKeys
	}

	for _, k := range keys {
		key, _, _, _, err := ssh.ParseAuthorizedKey([]byte(k))
		if err != nil {
			continue
		}
		a := key.Marshal()
		b := pubkey.Marshal()
		if len(a) == len(b) && subtle.ConstantTimeCompare(a, b) == 1 {
			return nil
		}
	}
	return ErrNoMatchingKeys
}
