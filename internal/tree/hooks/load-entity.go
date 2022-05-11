package hooks

import (
	"context"

	"google.golang.org/protobuf/proto"

	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

// LoadEntity loads an entity from the database.
type LoadEntity struct {
	tree.BaseHook
}

// Run attempts to load the entity specified by de.ID and if
// successful performs a deepcopy into the address pointed to by e.
// Any errors returned will be from the data storage layer.
func (l *LoadEntity) Run(ctx context.Context, e, de *pb.Entity) error {
	// This is a bit odd because we only get an address for e, not
	// the ability to point it somewhere else, so anything we want
	// to do that alters the initial contents needs to be copied
	// in.

	le, err := l.Storage().LoadEntity(ctx, de.GetID())
	if err != nil {
		return err
	}
	proto.Merge(e, le)

	return nil
}

func init() {
	startup.RegisterCallback(loadEntityCB)
}

func loadEntityCB() {
	tree.RegisterEntityHookConstructor("load-entity", NewLoadEntity)
}

// NewLoadEntity returns an initialized hook ready for use.
func NewLoadEntity(opts ...tree.HookOption) (tree.EntityHook, error) {
	opts = append([]tree.HookOption{
		tree.WithHookName("load-entity"),
		tree.WithHookPriority(0),
	}, opts...)

	return &LoadEntity{tree.NewBaseHook(opts...)}, nil
}
