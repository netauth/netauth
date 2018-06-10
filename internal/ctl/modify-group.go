package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/NetAuth/NetAuth/pkg/client"
	pb "github.com/NetAuth/Protocol"

	"github.com/google/subcommands"
)

type ModifyGroupCmd struct {
	name        string
	displayName string
	managedby   string
}

func (*ModifyGroupCmd) Name() string     { return "modify-group" }
func (*ModifyGroupCmd) Synopsis() string { return "Modify mutable fields on a group" }
func (*ModifyGroupCmd) Usage() string {
	return `modify-group --name <name> [fields-to-be-modified]
Modify a group by updating the named fields to the provided values.
`
}

func (p *ModifyGroupCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.name, "name", "", "Name of the group to modify")
	f.StringVar(&p.displayName, "display_name", "NO_CHANGE", "Group displayName")
	f.StringVar(&p.managedby, "managed_by", "NO_CHANGE", "Group that manages this group")
}

func (p *ModifyGroupCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
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

	group := &pb.Group{Name: &p.name}

	// This if block is kind of a hack, it is needed to ensure
	// that fields that weren't set to be modified in the command
	// line flags don't get set in the struct and so don't get
	// overwritten on the server.  If this isn't done here then it
	// means that the server only remembers the last thing to
	// change.  If of course you want to literally set a field to
	// "NO_CHANGE" this is somewhat impossible to do with the CLI.
	if p.displayName != "NO_CHANGE" {
		group.DisplayName = &p.displayName
	}

	if p.managedby != "NO_CHANGE" {
		group.ManagedBy = &p.managedby
	}

	result, err := c.ModifyGroupMeta(group, t)
	if result.GetMsg() != "" {
		fmt.Println(result.GetMsg())
	}
	if err != nil {
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
