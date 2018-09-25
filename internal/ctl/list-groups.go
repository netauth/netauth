package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
)

// ListGroupsCmd lists the groups for  either a specific entity, or of
// the entire server.
type ListGroupsCmd struct {
	fields    string
	entityID  string
	indirects bool
}

// Name of this cmdlet will be 'list-groups'
func (*ListGroupsCmd) Name() string { return "list-groups" }

// Synopsis returns the short form usage information.
func (*ListGroupsCmd) Synopsis() string { return "Obtain a list of groups" }

// Usage returns the long form usage information.
func (*ListGroupsCmd) Usage() string {
	return `list-groups --entity <ID> --indirects [--fields field1,field2,field3...]

This command will return a list of groups, additional formatting
options can be selected for additional information.  If an entity is
specified, then only groups on that entity will be returned.
`
}

// SetFlags sets the cmdlet specific flags
func (p *ListGroupsCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.fields, "fields", "", "Comma separated list of fields to display")
	f.StringVar(&p.entityID, "entity", "", "Entity to obtain groups for, blank for all groups")
	f.BoolVar(&p.indirects, "indirects", true, "Include indirect group memberships")
}

// Execute runs the cmdlet.
func (p *ListGroupsCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Grab a client
	c, err := getClient()
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	// Obtain group list
	list, err := c.ListGroups(p.entityID, p.indirects)
	if err != nil {
		return subcommands.ExitFailure
	}

	// Print the list
	for _, g := range list {
		printGroup(g, p.fields)
	}

	return subcommands.ExitSuccess
}
