package hooks

import (
	"sort"
	"strings"

	"github.com/netauth/netauth/internal/tree"
	"github.com/netauth/netauth/internal/tree/util"

	pb "github.com/netauth/protocol"
)

// PatchGroupExpansions handles the application of group expansions.
// Checks that verify the referential integrity of these expansions
// should have already run before executing this hook.
type PatchGroupExpansions struct {
	tree.BaseHook
}

// Run iterates through the expansions in dg and applies them to g.
// DROP expansions are processed with fuzzy group name matching.
func (*PatchGroupExpansions) Run(g, dg *pb.Group) error {
	exps := dg.GetExpansions()
	for i := range exps {
		parts := strings.SplitN(exps[i], ":", 2)
		if parts[0] == "INCLUDE" || parts[0] == "EXCLUDE" {
			g.Expansions = util.PatchStringSlice(g.Expansions, exps[i], true, true)
		} else if parts[0] == "DROP" {
			// Patch out with fuzzy matching, this will
			// innevitably come back to bite someone, but
			// until then its much faster than doing
			// strict checking on the group name.
			g.Expansions = util.PatchStringSlice(g.Expansions, parts[1], false, false)
		}
	}
	sort.Strings(g.Expansions)

	return nil
}

func init() {
	tree.RegisterGroupHookConstructor("patch-group-expansions", NewPatchGroupExpansions)
}

// NewPatchGroupExpansions returns an initialized hook for use.
func NewPatchGroupExpansions(tree.RefContext) (tree.GroupHook, error) {
	return &PatchGroupExpansions{tree.NewBaseHook("patch-group-expansions", 50)}, nil
}
