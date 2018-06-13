package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
)

type GroupExpansionsCmd struct {
	parent string
	child  string
	mode   string
}

func (*GroupExpansionsCmd) Name() string     { return "group-expansions" }
func (*GroupExpansionsCmd) Synopsis() string { return "Modify group expansions" }
func (*GroupExpansionsCmd) Usage() string {
	return `group-expansions --parent <parent> --child <child> --mode <INCLUDE|EXCLUDE|DROP>

Modify group expansions.  INCLUDE will include the children of the
named group in the parent, EXCLUDE will exclude the children of the
named group from the parent, and DROP will remove rules of either
type.`
}

func (p *GroupExpansionsCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.parent, "parent", "", "Parent Group")
	f.StringVar(&p.child, "child", "", "Child Group")
	f.StringVar(&p.mode, "mode", "INCLUDE", "Mode, must be one of INCLUDE, EXCLUDE, or DROP")
}

func (p *GroupExpansionsCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if p.parent == "" || p.child == "" {
		fmt.Println("--parent and --child must both be specified!")
		return subcommands.ExitFailure
	}

	// Grab a client
	c, err := getClient()
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	// Get the authorization token
	t, err := c.GetToken(entity, secret)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	result, err := c.ModifyGroupExpansions(t, p.parent, p.child, p.mode)
	if result.GetMsg() != "" {
		fmt.Println(result.GetMsg())
	}
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}
