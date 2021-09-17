package db

import (
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

type dummyKV struct{}

func newDummyKV(hclog.Logger) (KVStore, error)   { return &dummyKV{}, nil }
func (d *dummyKV) Put(string, []byte) error      { return nil }
func (d *dummyKV) Get(string) ([]byte, error)    { return nil, nil }
func (d *dummyKV) Del(string) error              { return nil }
func (d *dummyKV) Keys(string) ([]string, error) { return nil, nil }
func (d *dummyKV) Close() error                  { return nil }
func (d *dummyKV) Capabilities() []KVCapability  { return nil }
func (d *dummyKV) SetEventFunc(func(Event))      {}

func TestRegisterKV(t *testing.T) {
	kvBackends = make(map[string]KVFactory)
	RegisterKV("dummy", newDummyKV)
	assert.Len(t, kvBackends, 1)
	RegisterKV("dummy", newDummyKV)
	assert.Len(t, kvBackends, 1)
}

func TestNewKV(t *testing.T) {
	RegisterKV("dummy", newDummyKV)
	res, err := NewKV("dummy", hclog.NewNullLogger())
	assert.Nil(t, err)
	assert.Implements(t, new(KVStore), res)

	_, err = NewKV("does-not-exist", hclog.NewNullLogger())
	assert.Equal(t, err, ErrUnknownDatabase)
}
