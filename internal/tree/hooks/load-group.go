package hooks

import (
	"context"

	"google.golang.org/protobuf/proto"

	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

// LoadGroup loads an entity from the database.
type LoadGroup struct {
	tree.BaseHook
}

// Run attempts to load the group specified by de.Name and if
// successful performs a deepcopy into the address pointed to by g.
// Any errors returned will be from the data storage layer.
func (l *LoadGroup) Run(ctx context.Context, g, dg *pb.Group) error {
	// This is a bit odd because we only get an address for g, not
	// the ability to point it somewhere else, so anything we want
	// to do that alters the initial contents needs to be copied
	// in.

	lg, err := l.Storage().LoadGroup(ctx, dg.GetName())
	if err != nil {
		return err
	}
	proto.Merge(g, lg)

	return nil
}

func init() {
	startup.RegisterCallback(loadGroupCB)
}

func loadGroupCB() {
	tree.RegisterGroupHookConstructor("load-group", NewLoadGroup)
}

// NewLoadGroup returns an initialized hook ready for use.
func NewLoadGroup(opts ...tree.HookOption) (tree.GroupHook, error) {
	opts = append([]tree.HookOption{
		tree.WithHookName("load-group"),
		tree.WithHookPriority(0),
	}, opts...)

	return &LoadGroup{tree.NewBaseHook(opts...)}, nil
}
