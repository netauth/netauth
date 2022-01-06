package fs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/netauth/netauth/pkg/token/keyprovider"
)

func TestProvideNoKey(t *testing.T) {
	viper.Set("core.conf", t.TempDir())
	p, _ := newFS(hclog.NewNullLogger())

	_, err := p.Provide("local", "test")
	assert.Equal(t, err, keyprovider.ErrNoSuchKey)

}

func TestProvideBadPath(t *testing.T) {
	tbase := t.TempDir()
	viper.Set("core.conf", tbase)
	p, _ := newFS(hclog.NewNullLogger())

	path := filepath.Join(tbase, "keys", "local-test.tokenkey")
	assert.Nil(t, os.MkdirAll(path, 0755))
	_, err := p.Provide("local", "test")
	assert.Equal(t, err, keyprovider.ErrInternal)
}

func TestProvide(t *testing.T) {
	viper.Set("core.conf", t.TempDir())
	p, _ := newFS(hclog.NewNullLogger())

	keydir := filepath.Join(viper.GetString("core.conf"), "keys")
	assert.Nil(t, os.MkdirAll(keydir, 0755))
	assert.Nil(t, os.WriteFile(filepath.Join(keydir, "local-test.tokenkey"), []byte("Hello World!"), 0644))

	b, err := p.Provide("local", "test")
	assert.Nil(t, err)
	assert.Equal(t, b, []byte("Hello World!"))
}
