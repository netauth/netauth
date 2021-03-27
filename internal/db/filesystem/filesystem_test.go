package filesystem

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/netauth/netauth/internal/db"
)

func TestCB(t *testing.T) {
	cb()
}

func TestNewKV(t *testing.T) {
	res, err := newKV(hclog.NewNullLogger())
	assert.Nil(t, err)
	assert.Implements(t, new(db.KVStore), res)
}

func TestSetEventFunc(t *testing.T) {
	kv, err := newKV(hclog.NewNullLogger())
	assert.Nil(t, err)

	assert.Nil(t, kv.(*Filesystem).eF)
	kv.SetEventFunc(func(db.Event) {})
	assert.NotNil(t, kv.(*Filesystem).eF)
}

func TestPut(t *testing.T) {
	kv, err := newKV(hclog.NewNullLogger())
	assert.Nil(t, err)
	kv.(*Filesystem).basePath = t.TempDir()

	os.MkdirAll(filepath.Join(kv.(*Filesystem).basePath, "no-access-here"), 0000)
	f, err := os.Create(filepath.Join(kv.(*Filesystem).basePath, "collide-with-this"))
	assert.Nil(t, err)
	f.Close()

	assert.Nil(t, kv.Put("/foo/bar", []byte("bytes!")))
	assert.Equal(t, ErrPathEscape, kv.Put("../out/of/chroot", []byte("evil data")))
	assert.NotNil(t, kv.Put("no-access-here/key", []byte("some bytes")))
	assert.NotNil(t, kv.Put("collide-with-this/key", []byte("some bytes")))
}

func TestGet(t *testing.T) {
	kv, err := newKV(hclog.NewNullLogger())
	assert.Nil(t, err)
	kv.(*Filesystem).basePath = t.TempDir()
	os.MkdirAll(filepath.Join(kv.(*Filesystem).basePath, "not-a-file"), 0755)

	_, err = kv.Get("../out/of/chroot")
	assert.Equal(t, ErrPathEscape, err)

	assert.Nil(t, kv.Put("/data/foo", []byte("some bytes")))
	res, err := kv.Get("/data/foo")
	assert.Nil(t, err)
	assert.Equal(t, []byte("some bytes"), res)

	_, err = kv.Get("/not-a-file")
	assert.NotNil(t, err)

	_, err = kv.Get("/does-not-exist")
	assert.Equal(t, db.ErrNoValue, err)
}

func TestDel(t *testing.T) {
	kv, err := newKV(hclog.NewNullLogger())
	assert.Nil(t, err)
	kv.(*Filesystem).basePath = t.TempDir()
	os.MkdirAll(filepath.Join(kv.(*Filesystem).basePath, "nested", "directory"), 0755)

	assert.Nil(t, kv.Put("/data/foo", []byte("some bytes")))

	assert.Nil(t, kv.Del("/data/foo"))
	assert.Equal(t, ErrPathEscape, kv.Del("../out/of/chroot"))
	assert.Equal(t, db.ErrNoValue, kv.Del("/does/not/exist"))
	assert.NotNil(t, kv.Del("/nested"))
}

func TestKeys(t *testing.T) {
	kv, err := newKV(hclog.NewNullLogger())
	assert.Nil(t, err)
	kv.(*Filesystem).basePath = t.TempDir()

	assert.Nil(t, kv.Put("/foo/foo", []byte("some bytes")))
	assert.Nil(t, kv.Put("/bar/foo", []byte("some bytes")))
	assert.Nil(t, kv.Put("/baz/foo", []byte("some bytes")))

	res, err := kv.Keys("/*/*")
	assert.Nil(t, err)
	assert.Equal(t, []string{"/bar/foo", "/baz/foo", "/foo/foo"}, res) //lexical

	res, err = kv.Keys("/ba*/*")
	assert.Nil(t, err)
	assert.Equal(t, []string{"/bar/foo", "/baz/foo"}, res) //lexical

}

func TestClose(t *testing.T) {
	kv, err := newKV(hclog.NewNullLogger())
	assert.Nil(t, err)

	assert.Nil(t, kv.Close())
}

func TestCapabilities(t *testing.T) {
	kv, err := newKV(hclog.NewNullLogger())
	assert.Nil(t, err)
	kv.(*Filesystem).basePath = t.TempDir()

	assert.Equal(t, []db.KVCapability{}, kv.Capabilities())

	f, err := os.Create(filepath.Join(kv.(*Filesystem).basePath, ".mutable"))
	assert.Nil(t, err)
	f.Close()
	assert.Equal(t, []db.KVCapability{db.KVMutable}, kv.Capabilities())
}

type eventHandler struct{ mock.Mock }

func (eh *eventHandler) FireEvent(e db.Event) {
	eh.Called(e)
}

func TestFireEventForKey(t *testing.T) {
	ef := eventHandler{}

	kv, err := newKV(hclog.NewNullLogger())
	assert.Nil(t, err)

	kv.SetEventFunc(ef.FireEvent)

	ef.On("FireEvent", db.Event{PK: "entity1", Type: db.EventEntityUpdate})
	ef.On("FireEvent", db.Event{PK: "entity1", Type: db.EventEntityDestroy})
	ef.On("FireEvent", db.Event{PK: "group1", Type: db.EventGroupUpdate})
	ef.On("FireEvent", db.Event{PK: "group1", Type: db.EventGroupDestroy})

	kv.(*Filesystem).fireEventForKey("/entities/entity1", eventUpdate)
	ef.AssertCalled(t, "FireEvent", db.Event{PK: "entity1", Type: db.EventEntityUpdate})

	kv.(*Filesystem).fireEventForKey("/entities/entity1", eventDelete)
	ef.AssertCalled(t, "FireEvent", db.Event{PK: "entity1", Type: db.EventEntityDestroy})

	kv.(*Filesystem).fireEventForKey("/groups/group1", eventUpdate)
	ef.AssertCalled(t, "FireEvent", db.Event{PK: "group1", Type: db.EventGroupUpdate})

	kv.(*Filesystem).fireEventForKey("/groups/group1", eventDelete)
	ef.AssertCalled(t, "FireEvent", db.Event{PK: "group1", Type: db.EventGroupDestroy})

	kv.(*Filesystem).fireEventForKey("/not/an/event/key", eventUpdate)

	// This assertion validates that the call for a key outside
	// the fixed keyspace went to the default case and didn't
	// trigger an event.
	ef.AssertNumberOfCalls(t, "FireEvent", 4)
}
