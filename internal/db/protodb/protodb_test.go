package protodb

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/internal/db"

	pb "github.com/NetAuth/Protocol"
)

func mkTmpTestDir(t *testing.T) string {
	dir, err := ioutil.TempDir("/tmp", "pdbtest")
	if err != nil {
		t.Error(err)
	}

	dir = filepath.Join(dir, "pdb")

	return dir
}

func cleanTmpTestDir(dir string, t *testing.T) {
	// Remove the tmpdir, don't want to clutter the filesystem
	path, final := filepath.Split(dir)

	// Strip the added directory path if this came from
	// mkTmpTestDir
	if final == "pdb" {
		dir = path
	}

	if err := os.RemoveAll(dir); err != nil {
		t.Log(err)
	}
}

func TestDiscoverEntities(t *testing.T) {
	// This is a slight race condition since we're manipulating
	// flags, but this shouldn't actually be flaky.
	*dataRoot = mkTmpTestDir(t)
	defer cleanTmpTestDir(*dataRoot, t)
	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

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
}

func TestNextEntityNumber(t *testing.T) {
	// This is a slight race condition since we're manipulating
	// flags, but this shouldn't actually be flaky.
	*dataRoot = mkTmpTestDir(t)
	defer cleanTmpTestDir(*dataRoot, t)
	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	s := []struct {
		ID            string
		number        int32
		nextUIDNumber int32
	}{
		{"foo", 1, 2},
		{"bar", 2, 3},
		{"baz", 65, 66}, // Numbers may be missing in the middle
		{"fuu", 23, 66}, // Later additions shouldn't alter max
	}

	for _, c := range s {
		//  Make sure the entity actually gets added
		e := &pb.Entity{ID: proto.String(c.ID), Number: proto.Int32(c.number)}
		if err := x.SaveEntity(e); err != nil {
			t.Error(err)
		}

		// Validate that after a given mutation the number is
		// still what we expect it to be.
		next, err := x.NextEntityNumber()
		if err != nil {
			t.Error(err)
		}
		if next != c.nextUIDNumber {
			t.Errorf("Wrong next number; got: %v want %v", next, c.nextUIDNumber)
		}
	}
}

func TestDiscoverGroups(t *testing.T) {
	// This is a slight race condition since we're manipulating
	// flags, but this shouldn't actually be flaky.
	*dataRoot = mkTmpTestDir(t)
	defer cleanTmpTestDir(*dataRoot, t)
	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	// Discover and verify that there are no groups on disk
	l, err := x.DiscoverGroupNames()
	if err != nil {
		t.Error(err)
	}
	if len(l) != 0 {
		t.Error("DiscoverGroupNames made up a group")
	}

	// Write out a group
	if err := x.SaveGroup(&pb.Group{Name: proto.String("group1")}); err != nil {
		t.Error(err)
	}

	// Discover again and verify that there is now one group on
	// disk
	l, err = x.DiscoverGroupNames()
	if err != nil {
		t.Error(err)
	}
	if len(l) != 1 {
		t.Error("DiscoverGroupNames failed to discover the right number of groups")
	}
	if l[0] != "group1" {
		t.Errorf("DiscoverGroupNames discovered the wrong name: '%s'", l[0])
	}
}

func TestEntitySaveLoadDelete(t *testing.T) {
	// This is a slight race condition since we're manipulating
	// flags, but this shouldn't actually be flaky.
	*dataRoot = mkTmpTestDir(t)
	defer cleanTmpTestDir(*dataRoot, t)
	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

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
	if _, err := x.LoadEntity("foo"); err != db.ErrUnknownEntity {
		t.Error(err)
	}
}

func TestGroupSaveLoadDelete(t *testing.T) {
	// This is a slight race condition since we're manipulating
	// flags, but this shouldn't actually be flaky.
	*dataRoot = mkTmpTestDir(t)
	defer cleanTmpTestDir(*dataRoot, t)
	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	g := &pb.Group{Name: proto.String("foo")}

	// Write an entity to disk
	if err := x.SaveGroup(g); err != nil {
		t.Error(err)
	}

	// Load  it back, it  should still be  the same, but  we check
	// this to be sure.
	ng, err := x.LoadGroup("foo")
	if err != nil {
		t.Error(err)
	}
	if !proto.Equal(g, ng) {
		t.Errorf("Loaded group and original are not equivalent! '%v', '%v'", g, ng)
	}

	// Delete the group and confirm that loading it returns an
	// error.
	if err := x.DeleteGroup("foo"); err != nil {
		t.Error(err)
	}
	if _, err := x.LoadGroup("foo"); err != db.ErrUnknownGroup {
		t.Error(err)
	}
}

func TestEnsureDataDirectoryBadBase(t *testing.T) {
	// This is a slight race condition since we're manipulating
	// flags, but this shouldn't actually be flaky.
	*dataRoot = mkTmpTestDir(t)
	defer cleanTmpTestDir(*dataRoot, t)

	if _, err := os.OpenFile(*dataRoot, os.O_RDONLY|os.O_CREATE, 0000); err != nil {
		t.Fatal(err)
	}

	_, err := New()
	if err != db.ErrInternalError {
		t.Fatal(err)
	}
}

func TestEnsureDataDirectoryBadSubDir(t *testing.T) {
	// This is a slight race condition since we're manipulating
	// flags, but this shouldn't actually be flaky.
	*dataRoot = mkTmpTestDir(t)
	defer cleanTmpTestDir(*dataRoot, t)

	if err := os.Mkdir(*dataRoot, 0750); err != nil {
		t.Fatal(err)
	}

	if _, err := os.OpenFile(filepath.Join(*dataRoot, entitySubdir), os.O_RDONLY|os.O_CREATE, 0666); err != nil {
		t.Fatal(err)
	}
	_, err := New()
	if err != db.ErrInternalError {
		t.Fatal(err)
	}
}

func TestLoadEntityBadFile(t *testing.T) {
	// This is a slight race condition since we're manipulating
	// flags, but this shouldn't actually be flaky.
	*dataRoot = mkTmpTestDir(t)
	defer cleanTmpTestDir(*dataRoot, t)
	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	// Write an entity to disk
	e := &pb.Entity{ID: proto.String("foo")}
	if err := x.SaveEntity(e); err != nil {
		t.Error(err)
	}

	// Make the entity unreadable
	if err := os.Chmod(filepath.Join(*dataRoot, entitySubdir, "foo.dat"), 0000); err != nil {
		t.Error(err)
	}

	if _, err := x.LoadEntity("foo"); err != db.ErrInternalError {
		t.Error(err)
	}
}

func TestLoadEntityBadParse(t *testing.T) {
	// This is a slight race condition since we're manipulating
	// flags, but this shouldn't actually be flaky.
	*dataRoot = mkTmpTestDir(t)
	defer cleanTmpTestDir(*dataRoot, t)
	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	// Write an empty file to disk that will fail to unmarshal
	if err := ioutil.WriteFile(filepath.Join(*dataRoot, entitySubdir, "foo.dat"), []byte("foo"), 0666); err != nil {
		t.Fatal(err)
	}

	if _, err := x.LoadEntity("foo"); err != db.ErrInternalError {
		t.Error(err)
	}

}

func TestSaveEntityBadEntity(t *testing.T) {
	// This is a slight race condition since we're manipulating
	// flags, but this shouldn't actually be flaky.
	*dataRoot = mkTmpTestDir(t)
	defer cleanTmpTestDir(*dataRoot, t)
	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	if err := x.SaveEntity(nil); err != db.ErrInternalError {
		t.Error(err)
	}
}

func TestSaveEntityUnwritableFile(t *testing.T) {
	// This is a slight race condition since we're manipulating
	// flags, but this shouldn't actually be flaky.
	*dataRoot = mkTmpTestDir(t)
	defer cleanTmpTestDir(*dataRoot, t)
	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	// Write an empty file to disk that will collide with a proper write
	if err := ioutil.WriteFile(filepath.Join(*dataRoot, entitySubdir, "foo.dat"), []byte("foo"), 0000); err != nil {
		t.Fatal(err)
	}

	// Attempt to write an entity to disk
	e := &pb.Entity{ID: proto.String("foo")}
	if err := x.SaveEntity(e); err != db.ErrInternalError {
		t.Error(err)
	}
}

func TestDeleteUnknownEntity(t *testing.T) {
	// This is a slight race condition since we're manipulating
	// flags, but this shouldn't actually be flaky.
	*dataRoot = mkTmpTestDir(t)
	defer cleanTmpTestDir(*dataRoot, t)
	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	if err := x.DeleteEntity("unknown-entity"); err != db.ErrUnknownEntity {
		t.Error(err)
	}
}

func TestLoadGroupBadFile(t *testing.T) {
	// This is a slight race condition since we're manipulating
	// flags, but this shouldn't actually be flaky.
	*dataRoot = mkTmpTestDir(t)
	defer cleanTmpTestDir(*dataRoot, t)
	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	// Write a group to disk
	e := &pb.Group{Name: proto.String("group1")}
	if err := x.SaveGroup(e); err != nil {
		t.Error(err)
	}

	// Make the group unreadable
	if err := os.Chmod(filepath.Join(*dataRoot, groupSubdir, "group1.dat"), 0000); err != nil {
		t.Error(err)
	}

	if _, err := x.LoadGroup("group1"); err != db.ErrInternalError {
		t.Error(err)
	}
}

func TestLoadGroupBadParse(t *testing.T) {
	// This is a slight race condition since we're manipulating
	// flags, but this shouldn't actually be flaky.
	*dataRoot = mkTmpTestDir(t)
	defer cleanTmpTestDir(*dataRoot, t)
	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	// Write an empty file to disk that will fail to unmarshal
	if err := ioutil.WriteFile(filepath.Join(*dataRoot, groupSubdir, "group1.dat"), []byte("group1"), 0666); err != nil {
		t.Fatal(err)
	}

	if _, err := x.LoadGroup("group1"); err != db.ErrInternalError {
		t.Error(err)
	}

}

func TestSaveGroupBadGroup(t *testing.T) {
	// This is a slight race condition since we're manipulating
	// flags, but this shouldn't actually be flaky.
	*dataRoot = mkTmpTestDir(t)
	defer cleanTmpTestDir(*dataRoot, t)
	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	if err := x.SaveGroup(nil); err != db.ErrInternalError {
		t.Error(err)
	}
}

func TestSaveGroupUnwritableFile(t *testing.T) {
	// This is a slight race condition since we're manipulating
	// flags, but this shouldn't actually be flaky.
	*dataRoot = mkTmpTestDir(t)
	defer cleanTmpTestDir(*dataRoot, t)
	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	// Write an empty file to disk that will collide with a proper write
	if err := ioutil.WriteFile(filepath.Join(*dataRoot, groupSubdir, "group1.dat"), []byte("group1"), 0000); err != nil {
		t.Fatal(err)
	}

	// Attempt to write an entity to disk
	e := &pb.Group{Name: proto.String("group1")}
	if err := x.SaveGroup(e); err != db.ErrInternalError {
		t.Error(err)
	}
}

func TestDeleteUnknownGroup(t *testing.T) {
	// This is a slight race condition since we're manipulating
	// flags, but this shouldn't actually be flaky.
	*dataRoot = mkTmpTestDir(t)
	defer cleanTmpTestDir(*dataRoot, t)
	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	if err := x.DeleteGroup("unknown-group"); err != db.ErrUnknownGroup {
		t.Error(err)
	}
}

func TestNextGroupNumber(t *testing.T) {
	x, err := New()
	if err != nil {
		t.Fatal(err)
	}
	s := []struct {
		name   string
		number int32
		want   int32
	}{
		{"foo", 1, 2},
		{"bar", 2, 3},
		{"baz", 65, 66}, // Numbers may be missing in the middle
		{"fuu", 23, 66}, // Later additions shouldn't alter max
	}

	for _, c := range s {
		g := &pb.Group{Name: proto.String(c.name), Number: proto.Int32(c.number)}
		if err := x.SaveGroup(g); err != nil {
			t.Error(err)
		}

		// Validate that after a given mutation the number is
		// still what we expect it to be.
		next, err := x.NextGroupNumber()
		if err != nil {
			t.Error(err)
		}
		if next != c.want {
			t.Errorf("Wrong next number; got: %v want %v", next, c.want)
		}
	}
}

func TestHealthCheckOK(t *testing.T) {
	*dataRoot = mkTmpTestDir(t)
	defer cleanTmpTestDir(*dataRoot, t)
	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	rx, ok := x.(*ProtoDB)
	if !ok {
		t.Fatal("Bad type assertion")
	}

	status := rx.healthCheck()

	if !status.OK || status.Status != "ProtoDB is operating normally" {
		t.Errorf("Bad status: %v", status)
	}
}

func TestHealthCheckBadBase(t *testing.T) {
	*dataRoot = mkTmpTestDir(t)
	defer cleanTmpTestDir(*dataRoot, t)

	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	if err := os.RemoveAll(*dataRoot); err != nil {
		t.Fatal(err)
	}
	if _, err := os.OpenFile(*dataRoot, os.O_RDONLY|os.O_CREATE, 0666); err != nil {
		t.Fatal(err)
	}

	rx, ok := x.(*ProtoDB)
	if !ok {
		t.Fatal("Bad type assertion")
	}

	status := rx.healthCheck()

	if status.OK {
		t.Errorf("Bad status: %v", status)
	}
}

func TestHealthCheckBadPermissions(t *testing.T) {
	*dataRoot = mkTmpTestDir(t)
	defer cleanTmpTestDir(*dataRoot, t)

	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Chmod(*dataRoot, 0770); err != nil {
		t.Fatal(err)
	}

	rx, ok := x.(*ProtoDB)
	if !ok {
		t.Fatal("Bad type assertion")
	}

	status := rx.healthCheck()

	if status.OK {
		t.Errorf("Bad status: %v", status)
	}
}

func TestHealthCheckBadStat(t *testing.T) {
	*dataRoot = mkTmpTestDir(t)
	defer cleanTmpTestDir(*dataRoot, t)

	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	rx, ok := x.(*ProtoDB)
	if !ok {
		t.Fatal("Bad type assertion")
	}

	rx.dataRoot = filepath.Join(*dataRoot, "bad")

	status := rx.healthCheck()

	if status.OK {
		t.Errorf("Bad status: %v", status)
	}
}
