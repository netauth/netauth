// Package fs implements a filesystem based token cache.  This is the
// cache that most programs will want to use since the tokens can be
// pre-fetched by the system's login tasks or by the NetAuth CLI.
//
// This package has known data races and does not use atomic file
// writes because this functionality is hilariously broken on MS
// Windows, and the risk is tolerable.
package fs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/netauth/netauth/pkg/token/cache"
)

const (
	extension = "nt"
)

type fsCache struct {
	basepath string
}

func init() {
	cache.RegisterTokenCacheFactory("fs", new)
}

func new() (cache.TokenCache, error) {
	c := &fsCache{
		basepath: os.TempDir(),
	}
	return c, nil
}

// PutToken will write the token to a file in the specified basepath.
// The file will be written with the permissions of the curent user
// and will not be readable to other users.
func (fc *fsCache) PutToken(owner, token string) error {
	return ioutil.WriteFile(fc.filepathFromOwner(owner), []byte(token), 0600)
}

// GetToken will retrieve a file of the form <owner>.<extension> and
// return its contents.  The file must be readable to the user in the
// calling context.
func (fc *fsCache) GetToken(owner string) (string, error) {
	d, err := ioutil.ReadFile(fc.filepathFromOwner(owner))
	switch {
	case os.IsNotExist(err):
		return "", cache.ErrNoCachedToken
	case err != nil:
		return "", err
	}
	return string(d), nil
}

// DelToken will remove a token from the filesystem blindly.  It does
// not check the contents, and does not check the validity of the
// token, so call with care.  The token to be removed must be
// writeable to the current user.
func (fc *fsCache) DelToken(owner string) error {
	if err := os.Remove(fc.filepathFromOwner(owner)); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// filepathFromOwner is a convenience function which encapsulates the
// path selection logic for where tokens are stored.
func (fc *fsCache) filepathFromOwner(owner string) string {
	return filepath.Join(
		fc.basepath,
		fmt.Sprintf("%s.%s", owner, extension),
	)
}
