package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/NetAuth/NetAuth/pkg/client"

	"github.com/google/subcommands"
)

type EntityOutOfGroupCmd struct {
	modentityname string
	groupname     string
}

func (*EntityOutOfGroupCmd) Name() string { return "remove-entity-from-group" }
func (*EntityOutOfGroupCmd) Synopsis() string {
	return "Remove an existing entity from an existing group"
}
func (*EntityOutOfGroupCmd) Usage() string {
	return `remove-entity-from-group --ID <ID> --group <name>

Remove the entity identified by <ID> from the group named by <name>.
Both the entity and the group must already exist.
`
}

func (c *EntityOutOfGroupCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.modentityname, "ID", entity, "ID of the entity to remove from the group")
	f.StringVar(&c.groupname, "group", "", "Name of the group to remove from")
}

func (cmd *EntityOutOfGroupCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Ensure that the secret has been obtained to authorize this
	// command
	ensureSecret()

	// Grab a client
	c, err := client.New(serverAddr, serverPort, serviceID, clientID)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	msg, err := c.RemoveEntityFromGroup(entity, secret, cmd.modentityname, cmd.groupname)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}
	fmt.Println(msg)
	return subcommands.ExitSuccess
}
