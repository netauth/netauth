package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"

	pb "github.com/NetAuth/Protocol"
)

// EntityMembershipCmd modifies the direct group membership for an
// entity.
type EntityMembershipCmd struct {
	entityID  string
	groupName string
	add       bool
	drop      bool
}

// Name of this cmdlet will be 'entity-membership'
func (*EntityMembershipCmd) Name() string { return "entity-membership" }

// Synopsis returns the short-form usage information
func (*EntityMembershipCmd) Synopsis() string {
	return "Add or remove an existing entity to an existing group"
}

// Usage returns the long form usage information.
func (*EntityMembershipCmd) Usage() string {
	return `entity-membership --entity <ID> --group <name> --<add|remove>

Add or remove the named entity from the named group.  Both the entity
and the group must exist already.
`
}

// SetFlags sets the cmdlet specific flags.
func (cmd *EntityMembershipCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&cmd.entityID, "entity", getEntity(), "ID of the entity to add to the group")
	f.StringVar(&cmd.groupName, "group", "", "Name of the group to add to")
	f.BoolVar(&cmd.add, "add", false, "Add the specified membership")
	f.BoolVar(&cmd.drop, "drop", false, "Drop the specified membership")
}

// Execute runs the cmdlet.
func (cmd *EntityMembershipCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
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

	result := &pb.SimpleResult{}
	if cmd.add {
		result, err = c.AddEntityToGroup(t, cmd.groupName, cmd.entityID)
	} else if cmd.drop {
		result, err = c.RemoveEntityFromGroup(t, cmd.groupName, cmd.entityID)
	} else {
		fmt.Println("You must specify either --add or --drop for this command!")
		return subcommands.ExitFailure
	}
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}
	fmt.Println(result.GetMsg())
	return subcommands.ExitSuccess
}
