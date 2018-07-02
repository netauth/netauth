package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
)

// NewGroupCmd requests the server to provision a new group.
type NewGroupCmd struct {
	name        string
	displayName string
	gid         int
	managedBy   string
}

// Name of this cmdlet will be 'new-group'
func (*NewGroupCmd) Name() string { return "new-group" }

// Synopsis returns the short-form usage.
func (*NewGroupCmd) Synopsis() string { return "Add a new group to the server" }

// Usage returns the long format usage information.
func (*NewGroupCmd) Usage() string {
	return `new-group --name <name> [--display_name <display name>] [--gid_number <number>] [--managed_by <name>]
Allocate a new group with the given name and optional display name.
If the gid_number is not specified then the next available number will
be used.  The name and number cannot be changed once set, only the
displayName.
`
}

// SetFlags sets the cmdlet specific flags.
func (p *NewGroupCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.name, "name", "", "Name for the new group.")
	f.StringVar(&p.displayName, "display_name", "", "Display Name for the new group.")
	f.IntVar(&p.gid, "gid_number", -1, "Group ID Number for the new group (automatic if unset)")
	f.StringVar(&p.managedBy, "managed_by", "", "Group that will manage the new group")
}

// Execute runs the cmdlet.
func (p *NewGroupCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Grab a client
	c, err := getClient()
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	// Get the authorization token
	t, err := c.GetToken(getEntity(), getSecret())
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	result, err := c.NewGroup(p.name, p.displayName, p.managedBy, t, p.gid)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	fmt.Println(result.GetMsg())
	return subcommands.ExitSuccess
}
