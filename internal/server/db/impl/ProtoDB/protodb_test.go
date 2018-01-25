package ProtoDB

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/golang/protobuf/proto"

	pb "github.com/NetAuth/NetAuth/pkg/proto"
)

func TestDiscoverEntities(t *testing.T) {
	// This is a slight race condition since we're manipulating
	// flags, but this shouldn't actually be flaky.
	dir, err := ioutil.TempDir("/tmp", "pdbtest")
	if err != nil {
		t.Error(err)
	}
	*data_root = dir
	x := New()
	l, err := x.DiscoverEntityIDs()
	if err != nil {
		t.Error(err)
	}

	// At this point there are no entities, so the length should
	// be 0.
	if len(l) != 0 {
		t.Error("DiscoverEntityIDs made up an entity!")
	}

	// We'll save an entity that just has the ID set.  This isn't
	// very realistic, but its the minimum data needed to put a
	// file on disk.
	if err := x.SaveEntity(&pb.Entity{ID: proto.String("foo")}); err != nil {
		t.Error(err)
	}

	// Rerun discovery.
	l, err = x.DiscoverEntityIDs()
	if err != nil {
		t.Error(err)
	}

	// Now there should be one file on disk, and the ID that we've
	// discovered should be 'foo'
	if len(l) != 1 {
		t.Error("DiscoverEntityIDs failed to discover any entities!")
	}
	if l[0] != "foo" {
		t.Errorf("DiscoverEntityIDs discovered the wrong name: '%s'", l[0])
	}

	// Remove the tmpdir, don't want to clutter the filesystem
	if err := os.RemoveAll(dir); err != nil {
		t.Log(err)
	}
}

func TestSaveLoadDelete(t *testing.T) {
	// This is a slight race condition since we're manipulating
	// flags, but this shouldn't actually be flaky.
	dir, err := ioutil.TempDir("/tmp", "pdbtest")
	if err != nil {
		t.Error(err)
	}
	*data_root = dir
	x := New()

	e := &pb.Entity{ID: proto.String("foo")}

	// Write an entity to disk
	if err := x.SaveEntity(e); err != nil {
		t.Error(err)
	}

	// Load  it back, it  should still be  the same, but  we check
	// this to be sure.
	ne, err := x.LoadEntity("foo")
	if err != nil {
		t.Error(err)
	}
	if !proto.Equal(e, ne) {
		t.Errorf("Loaded entity and original are not equivalent! '%v', '%v'", e, ne)
	}

	// Delete the entity and confirm that loading it returns an
	// error.
	if err := x.DeleteEntity("foo"); err != nil {
		t.Error(err)
	}
	if _, err := x.LoadEntity("foo"); !os.IsNotExist(err) {
		t.Error(err)
	}

	// Remove the tmpdir, don't want to clutter the filesystem
	if err := os.RemoveAll(dir); err != nil {
		t.Log(err)
	}
}
