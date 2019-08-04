package protodb

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/radovskyb/watcher"
	"github.com/spf13/viper"

	"github.com/NetAuth/NetAuth/internal/db"

	pb "github.com/NetAuth/Protocol"
)

func mkTmpTestDir(t *testing.T) string {
	dir, err := ioutil.TempDir("/tmp", "pdbtest")
	if err != nil {
		t.Error(err)
	}

	return dir
}

func cleanTmpTestDir(dir string, t *testing.T) {
	db.DeregisterCallback("BleveIndexer")
	if err := os.RemoveAll(dir); err != nil {
		t.Log(err)
	}
}

func cleanUpWatcher(d db.DB, t *testing.T) {
	rx, ok := d.(*ProtoDB)
	if !ok {
		t.Fatal("Bad type assertion")
	}

	rx.w.Wait()
	rx.w.Close()
}

func TestDiscoverEntities(t *testing.T) {
	r := mkTmpTestDir(t)
	viper.Set("core.home", r)
	defer cleanTmpTestDir(r, t)

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
	r := mkTmpTestDir(t)
	viper.Set("core.home", r)
	defer cleanTmpTestDir(r, t)
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

func TestSearchEntities(t *testing.T) {
	r := mkTmpTestDir(t)
	viper.Set("core.home", r)
	defer cleanTmpTestDir(r, t)
	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	s := []struct {
		ID     string
		secret string
	}{
		{"entity1", "secret1"},
		{"entity2", "secret2"},
		{"entity3", "secret3"},
	}
	for i := range s {
		e := pb.Entity{ID: &s[i].ID, Secret: &s[i].secret}
		if err := x.SaveEntity(&e); err != nil {
			t.Fatal(err)
		}
	}

	res, err := x.SearchEntities(db.SearchRequest{})
	if err != db.ErrBadSearch {
		t.Fatal(err)
	}
	res, err = x.SearchEntities(db.SearchRequest{Expression: "entity1"})
	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 1 || res[0].GetID() != "entity1" {
		t.Log(res)
		t.Error("Result does not match expected singular value")
	}
}

func TestDiscoverGroups(t *testing.T) {
	r := mkTmpTestDir(t)
	viper.Set("core.home", r)
	defer cleanTmpTestDir(r, t)
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
	r := mkTmpTestDir(t)
	viper.Set("core.home", r)
	defer cleanTmpTestDir(r, t)
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
	r := mkTmpTestDir(t)
	viper.Set("core.home", r)
	defer cleanTmpTestDir(r, t)
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

func TestSearchGroups(t *testing.T) {
	r := mkTmpTestDir(t)
	viper.Set("core.home", r)
	defer cleanTmpTestDir(r, t)
	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	s := []struct {
		Name        string
		DisplayName string
	}{
		{"group1", "One"},
		{"group2", "Two"},
		{"group3", "Three"},
	}
	for i := range s {
		e := pb.Group{Name: &s[i].Name, DisplayName: &s[i].DisplayName}
		if err := x.SaveGroup(&e); err != nil {
			t.Fatal(err)
		}
	}

	res, err := x.SearchGroups(db.SearchRequest{})
	if err != db.ErrBadSearch {
		t.Fatal(err)
	}
	res, err = x.SearchGroups(db.SearchRequest{Expression: "group1"})
	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 1 || res[0].GetName() != "group1" {
		t.Log(res)
		t.Error("Result does not match expected singular value")
	}
}

func TestEnsureDataDirectoryBadBase(t *testing.T) {
	r := mkTmpTestDir(t)
	viper.Set("core.home", r)
	defer cleanTmpTestDir(r, t)

	if _, err := os.OpenFile(filepath.Join(r, "pdb"), os.O_RDONLY|os.O_CREATE, 0000); err != nil {
		t.Fatal(err)
	}

	_, err := New()
	if err != db.ErrInternalError {
		t.Fatal(err)
	}
}

func TestEnsureDataDirectoryBadSubDir(t *testing.T) {
	r := mkTmpTestDir(t)
	viper.Set("core.home", r)
	defer cleanTmpTestDir(r, t)

	if err := os.Mkdir(filepath.Join(r, "pdb"), 0750); err != nil {
		t.Fatal(err)
	}

	if _, err := os.OpenFile(filepath.Join(r, "pdb", entitySubdir), os.O_RDONLY|os.O_CREATE, 0666); err != nil {
		t.Fatal(err)
	}
	_, err := New()
	if err != db.ErrInternalError {
		t.Fatal(err)
	}
}

func TestLoadEntityBadFile(t *testing.T) {
	r := mkTmpTestDir(t)
	viper.Set("core.home", r)
	defer cleanTmpTestDir(r, t)
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
	if err := os.Chmod(filepath.Join(r, "pdb", entitySubdir, "foo.dat"), 0000); err != nil {
		t.Error(err)
	}

	if _, err := x.LoadEntity("foo"); err != db.ErrInternalError {
		t.Error(err)
	}
}

func TestLoadEntityBadParse(t *testing.T) {
	r := mkTmpTestDir(t)
	viper.Set("core.home", r)
	defer cleanTmpTestDir(r, t)
	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	// Write an empty file to disk that will fail to unmarshal
	if err := ioutil.WriteFile(filepath.Join(r, "pdb", entitySubdir, "foo.dat"), []byte("foo"), 0666); err != nil {
		t.Fatal(err)
	}

	if _, err := x.LoadEntity("foo"); err != db.ErrInternalError {
		t.Error(err)
	}

}

func TestSaveEntityBadEntity(t *testing.T) {
	r := mkTmpTestDir(t)
	viper.Set("core.home", r)
	defer cleanTmpTestDir(r, t)
	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	if err := x.SaveEntity(nil); err != db.ErrInternalError {
		t.Error(err)
	}
}

func TestSaveEntityUnwritableFile(t *testing.T) {
	r := mkTmpTestDir(t)
	viper.Set("core.home", r)
	defer cleanTmpTestDir(r, t)
	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	// Write an empty file to disk that will collide with a proper write
	if err := ioutil.WriteFile(filepath.Join(r, "pdb", entitySubdir, "foo.dat"), []byte("foo"), 0000); err != nil {
		t.Fatal(err)
	}

	// Attempt to write an entity to disk
	e := &pb.Entity{ID: proto.String("foo")}
	if err := x.SaveEntity(e); err != db.ErrInternalError {
		t.Error(err)
	}
}

func TestDeleteUnknownEntity(t *testing.T) {
	r := mkTmpTestDir(t)
	viper.Set("core.home", r)
	defer cleanTmpTestDir(r, t)
	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	if err := x.DeleteEntity("unknown-entity"); err != db.ErrUnknownEntity {
		t.Error(err)
	}
}

func TestLoadGroupBadFile(t *testing.T) {
	r := mkTmpTestDir(t)
	viper.Set("core.home", r)
	defer cleanTmpTestDir(r, t)
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
	if err := os.Chmod(filepath.Join(r, "pdb", groupSubdir, "group1.dat"), 0000); err != nil {
		t.Error(err)
	}

	if _, err := x.LoadGroup("group1"); err != db.ErrInternalError {
		t.Error(err)
	}
}

func TestLoadGroupBadParse(t *testing.T) {
	r := mkTmpTestDir(t)
	viper.Set("core.home", r)
	defer cleanTmpTestDir(r, t)
	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	// Write an empty file to disk that will fail to unmarshal
	if err := ioutil.WriteFile(filepath.Join(r, "pdb", groupSubdir, "group1.dat"), []byte("group1"), 0666); err != nil {
		t.Fatal(err)
	}

	if _, err := x.LoadGroup("group1"); err != db.ErrInternalError {
		t.Error(err)
	}

}

func TestSaveGroupBadGroup(t *testing.T) {
	r := mkTmpTestDir(t)
	viper.Set("core.home", r)
	defer cleanTmpTestDir(r, t)
	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	if err := x.SaveGroup(nil); err != db.ErrInternalError {
		t.Error(err)
	}
}

func TestSaveGroupUnwritableFile(t *testing.T) {
	r := mkTmpTestDir(t)
	viper.Set("core.home", r)
	defer cleanTmpTestDir(r, t)
	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	// Write an empty file to disk that will collide with a proper write
	if err := ioutil.WriteFile(filepath.Join(r, "pdb", groupSubdir, "group1.dat"), []byte("group1"), 0000); err != nil {
		t.Fatal(err)
	}

	// Attempt to write an entity to disk
	e := &pb.Group{Name: proto.String("group1")}
	if err := x.SaveGroup(e); err != db.ErrInternalError {
		t.Error(err)
	}
}

func TestDeleteUnknownGroup(t *testing.T) {
	r := mkTmpTestDir(t)
	viper.Set("core.home", r)
	defer cleanTmpTestDir(r, t)
	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	if err := x.DeleteGroup("unknown-group"); err != db.ErrUnknownGroup {
		t.Error(err)
	}
}

func TestNextGroupNumber(t *testing.T) {
	r := mkTmpTestDir(t)
	viper.Set("core.home", r)
	defer cleanTmpTestDir(r, t)
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
	r := mkTmpTestDir(t)
	viper.Set("core.home", r)
	defer cleanTmpTestDir(r, t)
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
	r := mkTmpTestDir(t)
	viper.Set("core.home", r)
	defer cleanTmpTestDir(r, t)

	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	if err := os.RemoveAll(r); err != nil {
		t.Fatal(err)
	}
	if _, err := os.OpenFile(r, os.O_RDONLY|os.O_CREATE, 0666); err != nil {
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

func TestHealthCheckNotDirectory(t *testing.T) {
	r := mkTmpTestDir(t)
	viper.Set("core.home", r)
	defer cleanTmpTestDir(r, t)

	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	eDir := filepath.Join(r, "pdb", entitySubdir)
	if err := os.RemoveAll(eDir); err != nil {
		t.Fatal(err)
	}
	if _, err := os.OpenFile(eDir, os.O_RDONLY|os.O_CREATE, 0666); err != nil {
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
	r := mkTmpTestDir(t)
	viper.Set("core.home", r)
	defer cleanTmpTestDir(r, t)

	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Chmod(filepath.Join(r, "pdb"), 0770); err != nil {
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
	r := mkTmpTestDir(t)
	viper.Set("core.home", r)
	defer cleanTmpTestDir(r, t)

	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	rx, ok := x.(*ProtoDB)
	if !ok {
		t.Fatal("Bad type assertion")
	}

	rx.dataRoot = filepath.Join(r, "bad")

	status := rx.healthCheck()

	if status.OK {
		t.Errorf("Bad status: %v", status)
	}
}

func TestIndexAvailableOnReload(t *testing.T) {
	r := mkTmpTestDir(t)
	defer cleanTmpTestDir(r, t)
	viper.Set("core.home", r)
	viper.Set("pdb.watcher", false)

	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	e1 := pb.Entity{ID: proto.String("entity1")}
	if err := x.SaveEntity(&e1); err != nil {
		t.Fatal(err)
	}

	g1 := pb.Group{Name: proto.String("group1")}
	if err := x.SaveGroup(&g1); err != nil {
		t.Fatal(err)
	}

	// After this point it is no longer safe to modify the DB
	// instance pointed to by X since callback hooks are global.
	db.DeregisterCallback("BleveIndexer")

	// Get a new pdb and make sure the index is populated
	y, err := New()
	if err != nil {
		t.Fatal(err)
	}

	res, err := y.SearchEntities(db.SearchRequest{Expression: "ID:entity1"})
	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 1 || res[0].GetID() != "entity1" {
		t.Log(res)
		t.Error("Result does not match expected singular value")
	}

	res2, err := y.SearchGroups(db.SearchRequest{Expression: "Name:group1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res2) != 1 || res2[0].GetName() != "group1" {
		t.Log(res2)
		t.Error("Result does not match expected singular value")
	}
}

func TestStartWatcher(t *testing.T) {
	r := mkTmpTestDir(t)
	viper.Set("core.home", r)
	viper.Set("pdb.watcher", true)
	viper.Set("pdb.watch-interval", "100ms")
	defer cleanTmpTestDir(r, t)

	x, err := New()
	defer cleanUpWatcher(x, t)
	if err != nil {
		t.Fatal(err)
	}
}

func TestWatcherLogging(t *testing.T) {
	r := mkTmpTestDir(t)
	viper.Set("core.home", r)
	viper.Set("pdb.watcher", true)
	viper.Set("pdb.watch-interval", "100ms")
	defer cleanTmpTestDir(r, t)

	x, err := New()
	defer cleanUpWatcher(x, t)
	if err != nil {
		t.Fatal(err)
	}

	rx, ok := x.(*ProtoDB)
	if !ok {
		t.Fatal("Bad type assertion")
	}

	rx.w.Error <- errors.New("Test Error")
}

func TestWatcherEvents(t *testing.T) {
	r := mkTmpTestDir(t)
	viper.Set("core.home", r)
	viper.Set("pdb.watcher", true)
	viper.Set("pdb.watch-interval", "100ms")
	defer cleanTmpTestDir(r, t)

	x, err := New()
	defer cleanUpWatcher(x, t)
	if err != nil {
		t.Fatal(err)
	}

	rx, ok := x.(*ProtoDB)
	if !ok {
		t.Fatal("Bad type assertion")
	}

	hookCorrect := false

	wf := func(e db.Event) {
		if e.PK == "entity1" && e.Type == db.EventEntityCreate {
			hookCorrect = true
		}
	}

	defer db.DeregisterCallback("TestWatcherEvents")
	db.RegisterCallback("TestWatcherEvents", wf)

	rx.w.Event <- watcher.Event{
		Op:   watcher.Create,
		Path: "pdb/entities/entity1.dat",
	}

	if !hookCorrect {
		t.Error("Watched event was incorrect")
	}
}

func TestConvertFSToDBEvent(t *testing.T) {
	cases := []struct {
		in   watcher.Event
		want db.Event
	}{
		{
			in: watcher.Event{
				Op:   watcher.Create,
				Path: "pdb/entities/entity1.dat",
			},
			want: db.Event{
				PK:   "entity1",
				Type: db.EventEntityCreate,
			},
		},
		{
			in: watcher.Event{
				Op:   watcher.Write,
				Path: "pdb/entities/entity1.dat",
			},
			want: db.Event{
				PK:   "entity1",
				Type: db.EventEntityUpdate,
			},
		},
		{
			in: watcher.Event{
				Op:   watcher.Remove,
				Path: "pdb/entities/entity1.dat",
			},
			want: db.Event{
				PK:   "entity1",
				Type: db.EventEntityDestroy,
			},
		},
		{
			in: watcher.Event{
				Op:   watcher.Create,
				Path: "pdb/groups/group1.dat",
			},
			want: db.Event{
				PK:   "group1",
				Type: db.EventGroupCreate,
			},
		},
		{
			in: watcher.Event{
				Op:   watcher.Write,
				Path: "pdb/groups/group1.dat",
			},
			want: db.Event{
				PK:   "group1",
				Type: db.EventGroupUpdate,
			},
		},
		{
			in: watcher.Event{
				Op:   watcher.Remove,
				Path: "pdb/groups/group1.dat",
			},
			want: db.Event{
				PK:   "group1",
				Type: db.EventGroupDestroy,
			},
		},
		{
			in: watcher.Event{
				Op:   watcher.Remove,
				Path: "pdb/unknown/group1.dat",
			},
			want: db.Event{},
		},
	}

	for i, c := range cases {
		if got := convertFSToDBEvent(c.in); got != c.want {
			t.Errorf("%d: Got %v Want %v", i, got, c.want)
		}
	}
}
