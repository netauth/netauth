package hooks

import (
	"context"
	"strings"

	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

// CheckExpansionTargets verifies that all expansions requested target
// groups that exist, unless the expansion type is DROP.
type CheckExpansionTargets struct {
	tree.BaseHook
}

// Run iterates through all expansions on dg and ensures that if the
// expansion type isn't DROP that the group actually exists.  This
// allows groups that have been deleted to effectively skip this
// check, since the only expansion that makes sense targeting a
// deleted group is to drop it.
func (cet *CheckExpansionTargets) Run(ctx context.Context, g, dg *pb.Group) error {
	targets := dg.GetExpansions()
	for i := range targets {
		parts := strings.SplitN(targets[i], ":", 2)
		if parts[0] == "DROP" {
			continue
		}
		if _, err := cet.Storage().LoadGroup(ctx, parts[1]); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	startup.RegisterCallback(checkExpansionTargetsCB)
}

func checkExpansionTargetsCB() {
	tree.RegisterGroupHookConstructor("check-expansion-targets", NewCheckExpansionTargets)
}

// NewCheckExpansionTargets returns a configured hook, ready for use.
func NewCheckExpansionTargets(opts ...tree.HookOption) (tree.GroupHook, error) {
	opts = append(
		[]tree.HookOption{
			tree.WithHookName("check-expansion-targets"),
			tree.WithHookPriority(40),
		}, opts...,
	)
	return &CheckExpansionTargets{tree.NewBaseHook(opts...)}, nil
}
