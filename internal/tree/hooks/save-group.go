package hooks

import (
	"github.com/NetAuth/NetAuth/internal/db"
	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

// SaveGroup is a hook intended to terminate processing chains by
// saving a modified group to the database.
type SaveGroup struct {
	tree.BaseHook
	db.DB
}

// Run will pass the group specified by g to the datastore and request
// it to be saved.
func (s *SaveGroup) Run(g, dg *pb.Group) error {
	return s.SaveGroup(g)
}

func init() {
	tree.RegisterGroupHookConstructor("save-group", NewSaveGroup)
}

// NewSaveGroup returns a configured hook for use.
func NewSaveGroup(c tree.RefContext) (tree.GroupProcessorHook, error) {
	return &SaveGroup{tree.NewBaseHook("save-group", 99), c.DB}, nil
}
