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
	"github.com/NetAuth/NetAuth/internal/server/entity_manager"
	"github.com/golang/protobuf/proto"

	pb "github.com/NetAuth/NetAuth/proto"
)

const entity_subdir = "entities"

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
func New() entity_manager.EMDiskInterface {
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
// by hand it should be safe enouth.
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

// LoadEntity loads a single entity fromt the data_root given the ID
// associated with the entity.
func (pdb *ProtoDB) LoadEntity(ID string) (*pb.Entity, error) {
	in, err := ioutil.ReadFile(filepath.Join(pdb.data_root, entity_subdir, fmt.Sprintf("%s.dat", ID)))
	if err != nil {
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

	if _, err := os.Stat(filepath.Join(pdb.data_root, "groups")); os.IsNotExist(err) {
		if err := os.Mkdir(filepath.Join(pdb.data_root, "groups"), 0755); err != nil {
			return err
		}
	}
	return nil
}
