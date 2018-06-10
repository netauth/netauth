package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/NetAuth/NetAuth/pkg/client"

	"github.com/google/subcommands"

	pb "github.com/NetAuth/Protocol"
)

type EntityMembershipCmd struct {
	entityID  string
	groupName string
	action    string
}

func (*EntityMembershipCmd) Name() string { return "entity-membership" }

func (*EntityMembershipCmd) Synopsis() string {
	return "Add or remove an existing entity to an existing group"
}

func (*EntityMembershipCmd) Usage() string {
	return `entity-membership --ID <ID> --group <name> --action <add|remove>

Add or remove the named entity from the named group.  Both the entity
and the group must exist already.
`
}

func (c *EntityMembershipCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.entityID, "ID", entity, "ID of the entity to add to the group")
	f.StringVar(&c.groupName, "group", "", "Name of the group to add to")
	f.StringVar(&c.action, "action", "", "Action to perform, must be 'add' or 'remove'")
}

func (cmd *EntityMembershipCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
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

	result := &pb.SimpleResult{}
	switch cmd.action {
	case "add":
		result, err = c.AddEntityToGroup(t, cmd.groupName, cmd.entityID)
	case "remove":
		result, err = c.RemoveEntityFromGroup(t, cmd.groupName, cmd.entityID)
	default:
		fmt.Println("You must specify either --action add or --action remove!")
		return subcommands.ExitFailure
	}
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}
	fmt.Println(result.GetMsg())
	return subcommands.ExitSuccess
}
