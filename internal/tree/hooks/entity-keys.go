package hooks

import (
	"strings"
	
	"github.com/NetAuth/NetAuth/internal/tree/util"

	pb "github.com/NetAuth/Protocol"
)

type AddEntityKey struct {}

func (*AddEntityKey) Name() string { return "add-entity-key" }
func (*AddEntityKey) Priority() int { return 50 }
func (*AddEntityKey) Run(e, de *pb.Entity) error {
	e.Meta.Keys = util.PatchStringSlice(e.Meta.Keys, de.Meta.Keys[0], true, true)
	return nil
}

type DelEntityKey struct {}

func (*DelEntityKey) Name() string { return "del-entity-key" }
func (*DelEntityKey) Priority() int { return 50 }
func (*DelEntityKey) Run(e, de *pb.Entity) error {
	key := strings.SplitN(de.Meta.Keys[0], ":", 2)
	e.Meta.Keys = util.PatchStringSlice(e.Meta.Keys, key[1], false, false)
	return nil
}
