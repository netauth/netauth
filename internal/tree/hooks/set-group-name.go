package hooks

import (
	pb "github.com/NetAuth/Protocol"
)

type SetGroupName struct {}

func (*SetGroupName) Name() string { return "set-group-name" }
func (*SetGroupName) Priority() int { return 50 }
func (*SetGroupName) Run(g, dg *pb.Group) error {
	g.Name = dg.Name
	return nil
}
