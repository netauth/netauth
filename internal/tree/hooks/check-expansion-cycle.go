package hooks

import (
	"context"
	"strings"

	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

// CheckExpansionCycles ensures that there is no path from the group g
// to expansions requested in the data group that could lead to a
// cycle in the inclusion graph.
type CheckExpansionCycles struct {
	tree.BaseHook
}

// Run will iterate through all expansions requested in dg and ensure
// that no cycles exist between the g and the requested include.  If
// the mode for any expansion is DROP that expansion will be skipped
// without checking.
func (cec *CheckExpansionCycles) Run(ctx context.Context, g, dg *pb.Group) error {
	exps := dg.GetExpansions()
	for i := range exps {
		parts := strings.SplitN(exps[i], ":", 2)
		// If the mode is DROP then it doesn't matter if it
		// conflicts.
		if parts[0] == "DROP" {
			continue
		}
		child, err := cec.Storage().LoadGroup(ctx, parts[1])
		if err != nil {
			return err
		}
		if cec.checkGroupCycles(ctx, child, g.GetName()) {
			return tree.ErrExistingExpansion
		}
	}
	return nil
}

// checkGroupCycles recurses down the group tree and tries to find the
// candidate group somewhere on the tree below the entry point.  The
// general usage would be to push in the target of the expansion as
// the group and then hunt for the parent group as the candidate.
func (cec *CheckExpansionCycles) checkGroupCycles(ctx context.Context, g *pb.Group, candidate string) bool {
	for _, exp := range g.GetExpansions() {
		parts := strings.SplitN(exp, ":", 2)
		if parts[1] == candidate {
			return true
		}
		ng, err := cec.Storage().LoadGroup(ctx, parts[1])
		if err != nil {
			// Play it safe, if we can't get the group
			// something may already be wrong.  Returning
			// true here can prevent further damage to the
			// tree.
			return true
		}
		if r := cec.checkGroupCycles(ctx, ng, candidate); r {
			return r
		}
	}
	return false
}

func init() { startup.RegisterCallback(checkExpansionCyclesCB) }

func checkExpansionCyclesCB() {
	tree.RegisterGroupHookConstructor("check-expansion-cycles", NewCheckExpansionCycles)
}

// NewCheckExpansionCycles returns a configured hook ready for use.
func NewCheckExpansionCycles(opts ...tree.HookOption) (tree.GroupHook, error) {
	opts = append([]tree.HookOption{
		tree.WithHookName("check-expansion-cycles"),
		tree.WithHookPriority(40),
	}, opts...)

	return &CheckExpansionCycles{tree.NewBaseHook(opts...)}, nil
}
