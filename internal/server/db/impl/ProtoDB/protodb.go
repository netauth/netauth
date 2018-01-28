package ProtoDB

// This is one of the simplest databases that just reads and writes
// protos to the local disk.  It's probably quite usable in
// environments that don't have high modification rates.

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/NetAuth/NetAuth/internal/server/db"
	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/pkg/errors"
	pb "github.com/NetAuth/NetAuth/pkg/proto"
)

const entity_subdir = "entities"
const group_subdir = "groups"

type ProtoDB struct {
	data_root string
}

var (
	data_root = flag.String("protodb_root", "./data", "Base directory for ProtoDB")
)

func init() {
	db.RegisterDB("ProtoDB", New)
}

// New returns a new ProtoDB instance that is initialized and ready
// for use.  This function will attempt to set up the data directory
// and fail out if it does not have permissions to write/stat the base
// directory and children.  This function will bail out the entire
// program as without the backing store the functionality of the rest
// of the server is undefined!
func New() db.EMDiskInterface {
	x := new(ProtoDB)
	x.data_root = *data_root
	if err := x.ensureDataDirectory(); err != nil {
		log.Fatalf("Could not establish data directory! (%s)", err)
		return nil
	}
	return x
}

// DiscoverEntityIDs returns a list of entity IDs that this loader can
// retrieve by globbing the entity directory of the data_root.  This
// is not foolproof, but assuming that the data_root is not modified
// by hand it should be safe enough.
func (pdb *ProtoDB) DiscoverEntityIDs() ([]string, error) {
	// Locate all known entities.
	globs, err := filepath.Glob(filepath.Join(pdb.data_root, entity_subdir, "*.dat"))
	if err != nil {
		return nil, err
	}

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
	in, err := ioutil.ReadFile(filepath.Join(pdb.data_root, entity_subdir, fmt.Sprintf("%s.dat", ID)))
	if err != nil {
		if os.IsNotExist(err) {
			// In the specific case of a non-existance,
			// that is a NO_ENTITY condition.
			return nil, errors.E_NO_ENTITY
		}
		log.Println("Error reading file:", err)
		return nil, err
	}
	e := &pb.Entity{}
	if err := proto.Unmarshal(in, e); err != nil {
		log.Printf("Failed to parse Entity from disk: (%s):", err)
		return nil, err
	}
	return e, nil
}

// LoadEntityNumber loads a single entity from the data_root given the
// uidNumber associated with the entity.
func (pdb *ProtoDB) LoadEntityNumber(number int32) (*pb.Entity, error) {
	l, err := pdb.DiscoverEntityIDs()
	if err != nil {
		return nil, err
	}

	for _, en := range l {
		e, err := pdb.LoadEntity(en)
		if err != nil {
			return nil, err
		}
		if e.GetUidNumber() == number {
			return e, nil
		}
	}
	return nil, errors.E_NO_ENTITY
}

// SaveEntity writes  an entity to  disk.  Errors may be  returned for
// proto marshal  errors or for  errors writing to disk.   No promises
// are made  regarding if  the data  has been written  to disk  at the
// return of this function as the operatig system may choose to buffer
// the data until a larger block may be written.
func (pdb *ProtoDB) SaveEntity(e *pb.Entity) error {
	out, err := proto.Marshal(e)
	if err != nil {
		log.Printf("Failed to marshal entity '%s' (%s)", e.GetID(), err)
	}

	if err := ioutil.WriteFile(filepath.Join(pdb.data_root, entity_subdir, fmt.Sprintf("%s.dat", e.GetID())), out, 0644); err != nil {
		log.Printf("Failed to aquire write handle for '%s'", e.GetID())
		return err
	}

	return nil
}

// DeleteEntity removes an entity from disk.  This is rather simple to
// do given that each entity is owned by exactly one file on disk.
// Simply removing the file is sufficient to delete the entity.
func (pdb *ProtoDB) DeleteEntity(ID string) error {
	return os.Remove(filepath.Join(pdb.data_root, entity_subdir, fmt.Sprintf("%s.dat", ID)))
}

// DiscoverGroupNames returns a list of entity IDs that this loader
// can retrieve by globbing the group directory of the data_root.
// This is not foolproof, but assuming that the data_root is not
// modified by hand it should be safe enough.
func (pdb *ProtoDB) DiscoverGroupNames() ([]string, error) {
	// Locate all known entities.
	globs, err := filepath.Glob(filepath.Join(pdb.data_root, group_subdir, "*.dat"))
	if err != nil {
		return nil, err
	}

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
	in, err := ioutil.ReadFile(filepath.Join(pdb.data_root, group_subdir, fmt.Sprintf("%s.dat", name)))
	if err != nil {
		if os.IsNotExist(err) {
			// This case is the group just flat not
			// existing and is returned as such.
			return nil, errors.E_NO_GROUP
		}
		log.Println("Error reading file:", err)
		return nil, err
	}
	e := &pb.Group{}
	if err := proto.Unmarshal(in, e); err != nil {
		log.Printf("Failed to parse Group from disk: (%s):", err)
		return nil, err
	}
	return e, nil
}

// LoadGroupNumber attempts to load a group by number.
func (pdb *ProtoDB) LoadGroupNumber(number int32) (*pb.Group, error) {
	l, err := pdb.DiscoverGroupNames()
	if err != nil {
		return nil, err
	}

	for _, gn := range l {
		g, err := pdb.LoadGroup(gn)
		if err != nil {
			return nil, err
		}
		if g.GetGidNumber() == number {
			return g, nil
		}
	}
	return nil, errors.E_NO_GROUP
}

// SaveGroup writes  an entity to  disk.  Errors may be  returned for
// proto marshal  errors or for  errors writing to disk.   No promises
// are made  regarding if  the data  has been written  to disk  at the
// return of this function as the operatig system may choose to buffer
// the data until a larger block may be written.
func (pdb *ProtoDB) SaveGroup(g *pb.Group) error {
	out, err := proto.Marshal(g)
	if err != nil {
		log.Printf("Failed to marshal entity '%s' (%s)", g.GetName(), err)
	}

	if err := ioutil.WriteFile(filepath.Join(pdb.data_root, group_subdir, fmt.Sprintf("%s.dat", g.GetName())), out, 0644); err != nil {
		log.Printf("Failed to aquire write handle for '%s'", g.GetName())
		return err
	}

	return nil
}

// DeleteGroup removes an entity from disk.  This is rather simple to
// do given that each group is owned by exactly one file on disk.
// Simply removing the file is sufficient to delete the entity.
func (pdb *ProtoDB) DeleteGroup(name string) error {
	err := os.Remove(filepath.Join(pdb.data_root, group_subdir, fmt.Sprintf("%s.dat", name)))

	if os.IsNotExist(err) {
		return errors.E_NO_GROUP
	}
	return err
}

// ensureDataDirectory is called during initialization of this backend
// to ensure that the data directories are available.
func (pdb *ProtoDB) ensureDataDirectory() error {
	if _, err := os.Stat(pdb.data_root); os.IsNotExist(err) {
		if err := os.Mkdir(pdb.data_root, 0755); err != nil {
			return err
		}
	}

	if _, err := os.Stat(filepath.Join(pdb.data_root, entity_subdir)); os.IsNotExist(err) {
		if err := os.Mkdir(filepath.Join(pdb.data_root, entity_subdir), 0755); err != nil {
			return err
		}
	}

	if _, err := os.Stat(filepath.Join(pdb.data_root, group_subdir)); os.IsNotExist(err) {
		if err := os.Mkdir(filepath.Join(pdb.data_root, group_subdir), 0755); err != nil {
			return err
		}
	}
	return nil
}
