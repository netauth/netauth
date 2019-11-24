package hooks

import (
	"strings"

	"github.com/netauth/netauth/internal/db"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

// CheckExpansionTargets verifies that all expansions requested target
// groups that exist, unless the expansion type is DROP.
type CheckExpansionTargets struct {
	tree.BaseHook
	db.DB
}

// Run iterates through all expansions on dg and ensures that if the
// expansion type isn't DROP that the group actually exists.  This
// allows groups that have been deleted to effectively skip this
// check, since the only expansion that makes sense targeting a
// deleted group is to drop it.
func (cet *CheckExpansionTargets) Run(g, dg *pb.Group) error {
	targets := dg.GetExpansions()
	for i := range targets {
		parts := strings.SplitN(targets[i], ":", 2)
		if parts[0] == "DROP" {
			continue
		}
		if _, err := cet.LoadGroup(parts[1]); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	tree.RegisterGroupHookConstructor("check-expansion-targets", NewCheckExpansionTargets)
}

// NewCheckExpansionTargets returns a configured hook, ready for use.
func NewCheckExpansionTargets(c tree.RefContext) (tree.GroupHook, error) {
	return &CheckExpansionTargets{tree.NewBaseHook("check-expansion-targets", 40), c.DB}, nil
}
