package fs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"

	"github.com/netauth/netauth/pkg/token/keyprovider"
)

func init() {
	keyprovider.Register("fs", newFS)
}

// FS retrieves key material from the filesystem.
type FS struct {
	l hclog.Logger
}

func newFS(l hclog.Logger) (keyprovider.KeyProvider, error) {
	return &FS{l: l.Named("fs")}, nil
}

// Provide returns key material.
func (fs FS) Provide(mech, usecase string) ([]byte, error) {
	b, err := os.ReadFile(
		filepath.Join(
			viper.GetString("core.conf"),
			"keys",
			fmt.Sprintf("%s-%s.tokenkey", mech, usecase),
		),
	)

	if os.IsNotExist(err) {
		return nil, keyprovider.ErrNoSuchKey
	}

	return b, nil
}
