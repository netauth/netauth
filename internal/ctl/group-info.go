package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
)

// GroupInfoCmd returns  information about a named  group filtered for
// specific fields.
type GroupInfoCmd struct {
	groupName   string
	fields string
}

// Name of this cmdlet will be 'group-info'
func (*GroupInfoCmd) Name() string { return "group-info" }

// Synopsis returns the short-form usage.
func (*GroupInfoCmd) Synopsis() string { return "Obtain information on a group" }

// Usage returns the long-form usage.
func (*GroupInfoCmd) Usage() string {
	return `group-info --group <name> [--fields field1,field2,field3...]

Return the fields of a group.  This will provide information on a
single group, as opposed to attempting to list all groups.
`
}

// SetFlags sets the cmdlet specific flags.
func (p *GroupInfoCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.fields, "fields", "", "Comma separated list of fields to display")
	f.StringVar(&p.groupName, "group", "", "Name of the group to query")
}

// Execute gets the group and prints information on it.
func (p *GroupInfoCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Grab a client
	c, err := getClient()
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	// Obtain group info
	result, err := c.GroupInfo(p.groupName)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	printGroup(result.GetGroup(), p.fields)

	if len(result.GetManaged()) > 0 {
		fmt.Printf("The following group(s) are managed by %s\n", p.groupName)
	}
	for _, gn := range result.GetManaged() {
		fmt.Printf("  - %s\n", gn)
	}
	return subcommands.ExitSuccess
}
