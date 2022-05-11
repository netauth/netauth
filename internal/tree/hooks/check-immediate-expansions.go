package hooks

import (
	"context"
	"strings"

	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

// CheckImmediateExpansions checks if a new expansion conflicts with
// an existing on attatched to the parent group. Expansions of type
// DROP are unchecked.
type CheckImmediateExpansions struct {
	tree.BaseHook
}

// Run iterates over the expansions already on g, for each comparing
// to each expansion in dg.  Excepting the case of an expansion type
// of DROP, which is unchecked, any matching expansion will result in
// an ErrExistingExpansion being returned.
func (cie *CheckImmediateExpansions) Run(_ context.Context, g, dg *pb.Group) error {
	existing := g.GetExpansions()
	proposed := dg.GetExpansions()
	for i := range proposed {
		parts := strings.SplitN(proposed[i], ":", 2)
		if parts[0] == "DROP" {
			continue
		}
		for k := range existing {
			if strings.Contains(existing[k], parts[1]) {
				return tree.ErrExistingExpansion
			}
		}
	}
	return nil
}

func init() {
	startup.RegisterCallback(checkImmediateExpansionsCB)
}

func checkImmediateExpansionsCB() {
	tree.RegisterGroupHookConstructor("check-immediate-expansions", NewCheckImmediateExpansions)
}

// NewCheckImmediateExpansions returns a configured hook for use.
func NewCheckImmediateExpansions(opts ...tree.HookOption) (tree.GroupHook, error) {
	opts = append(
		[]tree.HookOption{
			tree.WithHookName("check-immediate-expansions"),
			tree.WithHookPriority(40),
		}, opts...,
	)

	return &CheckImmediateExpansions{tree.NewBaseHook(opts...)}, nil
}
