// Package protodb is one of the simplest databases that just reads
// and writes protos to the local disk.  It's probably quite usable in
// environments that don't have high modification rates.
package protodb

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/golang/protobuf/proto"
	atomic "github.com/google/renameio"
	"github.com/hashicorp/go-hclog"
	"github.com/radovskyb/watcher"
	"github.com/spf13/viper"

	"github.com/netauth/netauth/internal/db"
	"github.com/netauth/netauth/internal/db/util"
	"github.com/netauth/netauth/internal/health"
	"github.com/netauth/netauth/internal/startup"

	pb "github.com/netauth/protocol"
)

const entitySubdir = "entities"
const groupSubdir = "groups"

// The ProtoDB type binds all methods that are a part of the protodb
// package.
type ProtoDB struct {
	dataRoot string
	idx      *util.SearchIndex

	w *watcher.Watcher
	l hclog.Logger
}

func init() {
	startup.RegisterCallback(cb)
}

func cb() {
	db.Register("ProtoDB", New)
}

// New returns a new ProtoDB instance that is initialized and ready
// for use.  This function will attempt to set up the data directory
// and fail out if it does not have permissions to write/stat the base
// directory and children.  This function will bail out the entire
// program as without the backing store the functionality of the rest
// of the server is undefined!
func New(l hclog.Logger) (db.DB, error) {
	x := new(ProtoDB)
	x.l = l.Named("protodb")
	x.dataRoot = filepath.Join(viper.GetString("core.home"), "pdb")
	x.idx = util.NewIndex()
	if err := x.ensureDataDirectory(); err != nil {
		x.l.Error("Could not establish data directory", "error", err)
		return nil, err
	}

	// ProtoDB uses callbacks to manage the search index, this
	// permits reactive handling of the index with a backgrounded
	// goroutine.
	x.idx.ConfigureCallback(x.LoadEntity, x.LoadGroup)
	db.RegisterCallback("BleveIndexer", x.idx.IndexCallback)

	// loadIndex triggers an index regenerate on server startup.
	x.loadIndex()

	// ProtoDB registers several health checks that allow the
	// system to know the status of the backend database.
	health.RegisterCheck("ProtoDB", x.healthCheck)

	if viper.GetBool("pdb.watcher") {
		x.l.Debug("Launching watcher")
		x.startWatcher()
	}

	return x, nil
}

// DiscoverEntityIDs returns a list of entity IDs that this loader can
// retrieve by globbing the entity directory of the data_root.  This
// is not foolproof, but assuming that the data_root is not modified
// by hand it should be safe enough.
func (pdb *ProtoDB) DiscoverEntityIDs() ([]string, error) {
	// Locate all known entities.  We throw away the error here
	// because from the manual: "Glob ignores file system errors
	// such as I/O errors reading directories. The only possible
	// returned error is ErrBadPattern, when pattern is
	// malformed."  Given that the pattern is hard coded, it is
	// impossible for an error to be returned from this call.
	globs, _ := filepath.Glob(filepath.Join(pdb.dataRoot, entitySubdir, "*.dat"))

	// Strip the extensions off the files.
	IDs := make([]string, 0)
	for _, g := range globs {
		f := filepath.Base(g)
		IDs = append(IDs, strings.Replace(f, ".dat", "", 1))
	}
	return IDs, nil
}

// LoadEntity loads a single entity from the data_root given the ID
// associated with the entity.
func (pdb *ProtoDB) LoadEntity(ID string) (*pb.Entity, error) {
	in, err := ioutil.ReadFile(filepath.Join(pdb.dataRoot, entitySubdir, fmt.Sprintf("%s.dat", ID)))
	if err != nil {
		if os.IsNotExist(err) {
			// In the specific case of a non-existence,
			// that is a UnknownEntity condition.
			return nil, db.ErrUnknownEntity
		}
		pdb.l.Error("Error reading file", "error", err)
		return nil, db.ErrInternalError
	}
	e := &pb.Entity{}
	if err := proto.Unmarshal(in, e); err != nil {
		pdb.l.Error("Error parsing entity", "error", err)
		return nil, db.ErrInternalError
	}
	return e, nil
}

// SaveEntity writes an entity to disk.  Errors may be returned for
// proto marshal errors or for errors writing to disk.  No promises
// are made regarding if the data has been written to disk at the
// return of this function as the operatig system may choose to buffer
// the data until a larger block may be written.
func (pdb *ProtoDB) SaveEntity(e *pb.Entity) error {
	out, err := proto.Marshal(e)
	if err != nil {
		pdb.l.Error("Failed to marshal entity", "entity", e.GetID(), "error", err)
		return db.ErrInternalError
	}

	eFile := filepath.Join(pdb.dataRoot, entitySubdir, fmt.Sprintf("%s.dat", e.GetID()))
	if err := atomic.WriteFile(eFile, out, 0644); err != nil {
		pdb.l.Error("Failed to write entity", "entity", e.GetID(), "error", err)
		return db.ErrInternalError
	}

	if !viper.GetBool("pdb.watcher") {
		pdb.l.Trace("Firing non-watched event", "event", db.EventEntityUpdate, "pk", e.GetID())
		db.FireEvent(db.Event{Type: db.EventEntityUpdate, PK: e.GetID()})
	}
	return nil
}

// DeleteEntity removes an entity from disk.  This is rather simple to
// do given that each entity is owned by exactly one file on disk.
// Simply removing the file is sufficient to delete the entity.
func (pdb *ProtoDB) DeleteEntity(ID string) error {
	err := os.Remove(filepath.Join(pdb.dataRoot, entitySubdir, fmt.Sprintf("%s.dat", ID)))

	if os.IsNotExist(err) {
		pdb.l.Warn("Attempt to delete unknown entity", "entity", ID)
		return db.ErrUnknownEntity
	} else if err != nil {
		pdb.l.Error("Error deleting entity", "entity", ID, "error", err)
		return db.ErrInternalError
	}

	if !viper.GetBool("pdb.watcher") {
		pdb.l.Trace("Firing non-watched event", "event", db.EventEntityDestroy, "pk", ID)
		db.FireEvent(db.Event{Type: db.EventEntityDestroy, PK: ID})
	}
	return nil
}

// NextEntityNumber computes and return the next entity number.
func (pdb *ProtoDB) NextEntityNumber() (int32, error) {
	return util.NextEntityNumber(pdb.LoadEntity, pdb.DiscoverEntityIDs)
}

// SearchEntities returns a slice of entity given a searchrequest.
func (pdb *ProtoDB) SearchEntities(r db.SearchRequest) ([]*pb.Entity, error) {
	res, err := pdb.idx.SearchEntities(r)
	if err != nil {
		return nil, err
	}
	return util.LoadEntityBatch(res, pdb.LoadEntity)
}

// DiscoverGroupNames returns a list of group names that this loader
// can retrieve by globbing the group directory of the data_root.
// This is not foolproof, but assuming that the data_root is not
// modified by hand it should be safe enough.
func (pdb *ProtoDB) DiscoverGroupNames() ([]string, error) {
	// Locate all known groups.  We throw away the error here
	// because from the manual: "Glob ignores file system errors
	// such as I/O errors reading directories. The only possible
	// returned error is ErrBadPattern, when pattern is
	// malformed."  Given that the pattern is hard coded, it is
	// impossible for an error to be returned from this call.
	globs, _ := filepath.Glob(filepath.Join(pdb.dataRoot, groupSubdir, "*.dat"))

	// Strip the extensions off the files.
	Names := make([]string, 0)
	for _, g := range globs {
		f := filepath.Base(g)
		Names = append(Names, strings.Replace(f, ".dat", "", 1))
	}
	return Names, nil
}

// LoadGroup attempts to load a group by name from the disk.  It can
// fail on proto errors or bogus file permissions reading the file.
func (pdb *ProtoDB) LoadGroup(name string) (*pb.Group, error) {
	in, err := ioutil.ReadFile(filepath.Join(pdb.dataRoot, groupSubdir, fmt.Sprintf("%s.dat", name)))
	if err != nil {
		if os.IsNotExist(err) {
			// This case is the group just flat not
			// existing and is returned as such.
			return nil, db.ErrUnknownGroup
		}
		pdb.l.Error("Error reading group", "group", name, "error", err)
		return nil, db.ErrInternalError
	}
	e := &pb.Group{}
	if err := proto.Unmarshal(in, e); err != nil {
		pdb.l.Error("Error parsing group", "group", name, "error", err)
		return nil, db.ErrInternalError
	}
	return e, nil
}

// SaveGroup writes an group to disk.  Errors may be returned for
// proto marshal errors or for errors writing to disk.  No promises
// are made regarding if the data has been written to disk at the
// return of this function as the operatig system may choose to buffer
// the data until a larger block may be written.
func (pdb *ProtoDB) SaveGroup(g *pb.Group) error {
	out, err := proto.Marshal(g)
	if err != nil {
		pdb.l.Error("Error marshaling group", "group", g.GetName(), "error", err)
		return db.ErrInternalError
	}

	gFile := filepath.Join(pdb.dataRoot, groupSubdir, fmt.Sprintf("%s.dat", g.GetName()))
	if err := atomic.WriteFile(gFile, out, 0644); err != nil {
		pdb.l.Error("Error writing group", "group", g.GetName(), "error", err)
		return db.ErrInternalError
	}

	if !viper.GetBool("pdb.watcher") {
		pdb.l.Trace("Firing non-watched event", "event", db.EventGroupUpdate, "pk", g.GetName())
		db.FireEvent(db.Event{Type: db.EventGroupUpdate, PK: g.GetName()})
	}
	return nil
}

// DeleteGroup removes a group from disk.  This is rather simple to do
// given that each group is owned by exactly one file on disk.  Simply
// removing the file is sufficient to delete the entity.
func (pdb *ProtoDB) DeleteGroup(name string) error {
	err := os.Remove(filepath.Join(pdb.dataRoot, groupSubdir, fmt.Sprintf("%s.dat", name)))

	if os.IsNotExist(err) {
		pdb.l.Warn("Attempt to remove non-existent group", "group", name)
		return db.ErrUnknownGroup
	} else if err != nil {
		pdb.l.Error("Error deleting group", "group", name, "error", err)
		return db.ErrInternalError
	}

	if !viper.GetBool("pdb.watcher") {
		pdb.l.Trace("Firing non-watched event", "event", db.EventGroupDestroy, "pk", name)
		db.FireEvent(db.Event{Type: db.EventGroupDestroy, PK: name})
	}
	return nil
}

// NextGroupNumber computes the next available group number.  This is
// very inefficient but it only is called when a new group is being
// created, which is hopefully infrequent.
func (pdb *ProtoDB) NextGroupNumber() (int32, error) {
	return util.NextGroupNumber(pdb.LoadGroup, pdb.DiscoverGroupNames)
}

// SearchGroups returns a slice of entity given a searchrequest.
func (pdb *ProtoDB) SearchGroups(r db.SearchRequest) ([]*pb.Group, error) {
	res, err := pdb.idx.SearchGroups(r)
	if err != nil {
		return nil, err
	}
	return util.LoadGroupBatch(res, pdb.LoadGroup)
}

// ensureDataDirectory is called during initialization of this backend
// to ensure that the data directories are available.
func (pdb *ProtoDB) ensureDataDirectory() error {
	dirs := []string{
		pdb.dataRoot,
		filepath.Join(pdb.dataRoot, entitySubdir),
		filepath.Join(pdb.dataRoot, groupSubdir),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0750); err != nil {
			pdb.l.Error("Error creating directory", "directory", d, "error", err)
			return db.ErrInternalError
		}
	}
	return nil
}

// healthCheck provides a sanity check that directories are okay and
// permissions are correct.
func (pdb *ProtoDB) healthCheck() health.SubsystemStatus {
	status := health.SubsystemStatus{
		OK:     true,
		Name:   "ProtoDB",
		Status: "ProtoDB is operating normally",
	}

	dirs := []string{
		pdb.dataRoot,
		filepath.Join(pdb.dataRoot, entitySubdir),
		filepath.Join(pdb.dataRoot, groupSubdir),
	}

	for _, dir := range dirs {
		stat, err := os.Stat(dir)
		if err != nil {
			status.OK = false
			status.Status = fmt.Sprintf("Error while checking '%s': '%s'", dir, err)
			return status
		}
		if !stat.IsDir() {
			status.OK = false
			status.Status = fmt.Sprintf("Path '%s' is not a directory", dir)
			return status
		}
		if stat.Mode().String() != "drwxr-x---" {
			status.OK = false
			status.Status = fmt.Sprintf("Path '%s' has the wrong permissions '%v'", dir, stat.Mode())
			return status
		}
	}
	return status
}

// loadIndex is used to fire an event for all entities and groups on
// the server to trigger an index action.
func (pdb *ProtoDB) loadIndex() {
	pdb.l.Debug("Beginning index regeneration")
	eList, _ := pdb.DiscoverEntityIDs()
	for i := range eList {
		pdb.l.Trace("Firing Index Event", "event", db.EventEntityUpdate, "pk", eList[i])
		db.FireEvent(db.Event{Type: db.EventEntityUpdate, PK: eList[i]})
	}

	gList, _ := pdb.DiscoverGroupNames()
	for i := range gList {
		pdb.l.Trace("Firing Index Event", "event", db.EventGroupUpdate, "pk", gList[i])
		db.FireEvent(db.Event{Type: db.EventGroupUpdate, PK: gList[i]})
	}
	pdb.l.Debug("Index regenerated")
}

// startWatcher is used to configure the filesystem watcher during
// startup and sets the directories to be watched, after
// configuration, it enables the watcher.
func (pdb *ProtoDB) startWatcher() {
	pdb.w = watcher.New()

	// We're only interested in events that would require
	// re-indexing.
	pdb.w.FilterOps(watcher.Create, watcher.Write, watcher.Remove)

	// As the directories themselves won't affect reindexing, we
	// can also just filter on dat files.
	r := regexp.MustCompile("^.*.dat$")
	pdb.w.AddFilterHook(watcher.RegexFilterHook(r, false))

	// While NetAuth itself won't write hidden files, its
	// plausible that a synchronization system might, so to be on
	// the safe side we ignore hidden files.
	pdb.w.IgnoreHiddenFiles(true)

	// Watch everything under the dataroot, which will include the
	// two directories for entities and groups, but we've filtered
	// out directories from the listing of changes that are
	// subscribed to.
	pdb.w.AddRecursive(pdb.dataRoot)

	go pdb.w.Start(viper.GetDuration("pdb.watch-interval"))
	go pdb.doWatch()
}

// doWatch is used to fire events if the watcher is enabled.
func (pdb *ProtoDB) doWatch() {
	for {
		select {
		case event := <-pdb.w.Event:
			e := pdb.convertFSToDBEvent(event)
			if !e.IsEmpty() {
				db.FireEvent(e)
			}
		case err := <-pdb.w.Error:
			pdb.l.Error("A watcher error has occurred", "error", err)
		case <-pdb.w.Closed:
			return
		}
	}
}

// convertFSToDBEvent figures out how to get from an event that
// converts from events happening on the filesystem to events that the
// database system understands.
func (pdb *ProtoDB) convertFSToDBEvent(e watcher.Event) db.Event {
	basename := filepath.Base(e.Path)

	ev := db.Event{
		PK: strings.TrimSuffix(basename, path.Ext(basename)),
	}

	subdir := filepath.Base(filepath.Dir(e.Path))
	if e.Op == watcher.Create && subdir == entitySubdir {
		ev.Type = db.EventEntityCreate
	} else if e.Op == watcher.Create && subdir == groupSubdir {
		ev.Type = db.EventGroupCreate
	} else if e.Op == watcher.Write && subdir == entitySubdir {
		ev.Type = db.EventEntityUpdate
	} else if e.Op == watcher.Write && subdir == groupSubdir {
		ev.Type = db.EventGroupUpdate
	} else if e.Op == watcher.Remove && subdir == entitySubdir {
		ev.Type = db.EventEntityDestroy
	} else if e.Op == watcher.Remove && subdir == groupSubdir {
		ev.Type = db.EventGroupDestroy
	} else {
		pdb.l.Warn("PDB Unmatched event!", "event", e)
		return db.Event{}
	}

	return ev
}
