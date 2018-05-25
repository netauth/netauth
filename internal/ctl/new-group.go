package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/NetAuth/NetAuth/pkg/client"

	"github.com/google/subcommands"
)

type NewGroupCmd struct {
	name        string
	displayName string
	gid         int
	managedBy   string
}

func (*NewGroupCmd) Name() string     { return "new-group" }
func (*NewGroupCmd) Synopsis() string { return "Add a new group to the server" }
func (*NewGroupCmd) Usage() string {
	return `new-group --name <name> [--display_name <display name>] [--gid_number <number>] [--managed_by <name>]
Allocate a new group with the given name and optional display name.
If the gid_number is not specified then the next available number will
be used.  The name and number cannot be changed once set, only the
displayName.
`
}

func (p *NewGroupCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.name, "name", "", "Name for the new group.")
	f.StringVar(&p.displayName, "display_name", "", "Display Name for the new group.")
	f.IntVar(&p.gid, "gid_number", -1, "Group ID Number for the new group (automatic if unset)")
	f.StringVar(&p.managedBy, "managed_by", "", "Group that will manage the new group")
}

func (p *NewGroupCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Grab a client
	c, err := client.New(serverAddr, serverPort, serviceID, clientID)
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

	msg, err := c.NewGroup(p.name, p.displayName, p.managedBy, t, p.gid)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	fmt.Println(msg)
	return subcommands.ExitSuccess
}
