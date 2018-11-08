package hooks

import (
	pb "github.com/NetAuth/Protocol"
)

type SetGroupDisplayName struct {}

func (*SetGroupDisplayName) Name() string { return "set-group-display-name" }
func (*SetGroupDisplayName) Priority() int { return 50 }
func (*SetGroupDisplayName) Run(g, dg *pb.Group) error {
	g.DisplayName = dg.DisplayName
	return nil
}
