package ctl

import (
	"context"
	"flag"
	"fmt"

	pb "github.com/NetAuth/Protocol"

	"github.com/google/subcommands"
)

// ModifyGroupCmd modifies mutable information on a group.
type ModifyGroupCmd struct {
	groupName   string
	displayName string
	managedby   string
}

// Name of this cmdlet is 'modify-group'
func (*ModifyGroupCmd) Name() string { return "modify-group" }

// Synopsis returns the short-form usage information.
func (*ModifyGroupCmd) Synopsis() string { return "Modify mutable fields on a group" }

// Usage returns the long-form usage information.
func (*ModifyGroupCmd) Usage() string {
	return `modify-group --group <name> [fields-to-be-modified]
Modify a group by updating the named fields to the provided values.
`
}

// SetFlags sets the cmdlet specific flags
func (p *ModifyGroupCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.groupName, "group", "", "Name of the group to modify")
	f.StringVar(&p.displayName, "display_name", "NO_CHANGE", "Group displayName")
	f.StringVar(&p.managedby, "managed_by", "NO_CHANGE", "Group that manages this group")
}

// Execute runs the cmdlet.
func (p *ModifyGroupCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
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

	group := &pb.Group{Name: &p.groupName}

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
