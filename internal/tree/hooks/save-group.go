package hooks

import (
	"github.com/netauth/netauth/internal/db"
	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
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
	startup.RegisterCallback(saveGroupCB)
}

func saveGroupCB() {
	tree.RegisterGroupHookConstructor("save-group", NewSaveGroup)
}

// NewSaveGroup returns a configured hook for use.
func NewSaveGroup(c tree.RefContext) (tree.GroupHook, error) {
	return &SaveGroup{tree.NewBaseHook("save-group", 99), c.DB}, nil
}
