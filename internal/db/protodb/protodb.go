// Package protodb is one of the simplest databases that just reads
// and writes protos to the local disk.  It's probably quite usable in
// environments that don't have high modification rates.
package protodb

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/NetAuth/NetAuth/internal/db"
	"github.com/NetAuth/NetAuth/internal/db/util"
	"github.com/NetAuth/NetAuth/internal/health"
	"github.com/golang/protobuf/proto"

	pb "github.com/NetAuth/Protocol"
)

const entitySubdir = "entities"
const groupSubdir = "groups"

// The ProtoDB type binds all methods that are a part of the protodb
// package.
type ProtoDB struct {
	dataRoot string
}

var (
	dataRoot = flag.String("protodb_root", "./data", "Base directory for ProtoDB")
)

func init() {
	db.Register("ProtoDB", New)
}

// New returns a new ProtoDB instance that is initialized and ready
// for use.  This function will attempt to set up the data directory
// and fail out if it does not have permissions to write/stat the base
// directory and children.  This function will bail out the entire
// program as without the backing store the functionality of the rest
// of the server is undefined!
func New() (db.DB, error) {
	x := new(ProtoDB)
	x.dataRoot = *dataRoot
	if err := x.ensureDataDirectory(); err != nil {
		log.Printf("Could not establish data directory! (%s)", err)
		return nil, err
	}

	health.RegisterCheck("ProtoDB", x.healthCheck)

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
		log.Println("Error reading file:", err)
		return nil, db.ErrInternalError
	}
	e := &pb.Entity{}
	if err := proto.Unmarshal(in, e); err != nil {
		log.Printf("Failed to parse Entity from disk: (%s):", err)
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
		log.Printf("Failed to marshal entity '%s' (%s)", e.GetID(), err)
		return db.ErrInternalError
	}

	if err := ioutil.WriteFile(filepath.Join(pdb.dataRoot, entitySubdir,
		fmt.Sprintf("%s.dat", e.GetID())), out, 0644); err != nil {
		log.Printf("Failed to acquire write handle for '%s'", e.GetID())
		return db.ErrInternalError
	}

	return nil
}

// DeleteEntity removes an entity from disk.  This is rather simple to
// do given that each entity is owned by exactly one file on disk.
// Simply removing the file is sufficient to delete the entity.
func (pdb *ProtoDB) DeleteEntity(ID string) error {
	err := os.Remove(filepath.Join(pdb.dataRoot, entitySubdir, fmt.Sprintf("%s.dat", ID)))

	if os.IsNotExist(err) {
		return db.ErrUnknownEntity
	}

	return nil
}

// NextEntityNumber computes and return the next entity number.
func (pdb *ProtoDB) NextEntityNumber() (int32, error) {
	return util.NextEntityNumber(pdb.LoadEntity, pdb.DiscoverEntityIDs)
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
		log.Println("Error reading file:", err)
		return nil, db.ErrInternalError
	}
	e := &pb.Group{}
	if err := proto.Unmarshal(in, e); err != nil {
		log.Printf("Failed to parse Group from disk: (%s):", err)
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
		log.Printf("Failed to marshal entity '%s' (%s)", g.GetName(), err)
		return db.ErrInternalError
	}

	if err := ioutil.WriteFile(filepath.Join(pdb.dataRoot, groupSubdir,
		fmt.Sprintf("%s.dat", g.GetName())), out, 0644); err != nil {
		log.Printf("Failed to acquire write handle for '%s'", g.GetName())
		return db.ErrInternalError
	}

	return nil
}

// DeleteGroup removes a group from disk.  This is rather simple to do
// given that each group is owned by exactly one file on disk.  Simply
// removing the file is sufficient to delete the entity.
func (pdb *ProtoDB) DeleteGroup(name string) error {
	err := os.Remove(filepath.Join(pdb.dataRoot, groupSubdir, fmt.Sprintf("%s.dat", name)))

	if os.IsNotExist(err) {
		return db.ErrUnknownGroup
	}

	return nil
}

// NextGroupNumber computes the next available group number.  This is
// very inefficient but it only is called when a new group is being
// created, which is hopefully infrequent.
func (pdb *ProtoDB) NextGroupNumber() (int32, error) {
	return util.NextGroupNumber(pdb.LoadGroup, pdb.DiscoverGroupNames)
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
			log.Println(err)
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
