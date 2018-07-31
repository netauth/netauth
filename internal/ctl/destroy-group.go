package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
)

// DestroyGroupCmd deletes a group
type DestroyGroupCmd struct {
	groupName string
}

// Name returns the name of this cmdlet.
func (*DestroyGroupCmd) Name() string { return "destroy-group" }

// Synopsis returns the short-form info for this cmdlet.
func (*DestroyGroupCmd) Synopsis() string { return "Delete a group existing on the server." }

// Usage returns the long-form info form this cmdlet.
func (*DestroyGroupCmd) Usage() string {
	return `destroy-group --group <name>
Delete the named group.
`
}

// SetFlags is the interface function which sets flags specific to this cmdlet.
func (p *DestroyGroupCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.groupName, "group", "", "Name of the group to destroy.")
}

// Execute is the interface function which runs this cmdlet.
func (p *DestroyGroupCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Grab a client
	c, err := getClient()
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	// Get the authorization token
	t, err := getToken(c, getEntity())
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	result, err := c.DeleteGroup(p.groupName, t)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	fmt.Println(result.GetMsg())
	return subcommands.ExitSuccess
}
