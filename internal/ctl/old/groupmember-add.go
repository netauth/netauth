package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/NetAuth/NetAuth/pkg/client"

	"github.com/google/subcommands"
)

type EntityIntoGroupCmd struct {
	modentityname string
	groupname     string
}

func (*EntityIntoGroupCmd) Name() string { return "add-entity-to-group" }
func (*EntityIntoGroupCmd) Synopsis() string { return "Add an existing entity to an existing group" }
func (*EntityIntoGroupCmd) Usage() string {
	return `add-entity-to-group --ID <ID> --group <name>

Add the entity identified by <ID> to the group named by <name>.  Both
the entity and the group must already exist.
`
}

func (c *EntityIntoGroupCmd) SetFlags( f*flag.FlagSet) {
	f.StringVar(&c.modentityname, "ID", entity, "ID of the entity to add to the group")
	f.StringVar(&c.groupname, "group", "", "Name of the group to add to")
}

func (cmd *EntityIntoGroupCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Ensure that the secret has been obtained to authorize this
	// command
	ensureSecret()

	// Grab a client
	c, err := client.New(serverAddr, serverPort, serviceID, clientID)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	msg, err := c.AddEntityToGroup(entity, secret, cmd.modentityname, cmd.groupname)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}
	fmt.Println(msg)
	return subcommands.ExitSuccess
}
