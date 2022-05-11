package hooks

import (
	"context"

	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

// SaveGroup is a hook intended to terminate processing chains by
// saving a modified group to the database.
type SaveGroup struct {
	tree.BaseHook
}

// Run will pass the group specified by g to the datastore and request
// it to be saved.
func (s *SaveGroup) Run(ctx context.Context, g, dg *pb.Group) error {
	return s.Storage().SaveGroup(ctx, g)
}

func init() {
	startup.RegisterCallback(saveGroupCB)
}

func saveGroupCB() {
	tree.RegisterGroupHookConstructor("save-group", NewSaveGroup)
}

// NewSaveGroup returns a configured hook for use.
func NewSaveGroup(opts ...tree.HookOption) (tree.GroupHook, error) {
	opts = append([]tree.HookOption{
		tree.WithHookName("save-group"),
		tree.WithHookPriority(99),
	}, opts...)

	return &SaveGroup{tree.NewBaseHook(opts...)}, nil
}
