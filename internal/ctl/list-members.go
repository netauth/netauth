package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
)

type ListMembersCmd struct {
	ID     string
	fields string
}

func (*ListMembersCmd) Name() string     { return "list-members" }
func (*ListMembersCmd) Synopsis() string { return "List members in a named group" }
func (*ListMembersCmd) Usage() string {
	return `list-members --group <group> [--fields field1,field2...]

List the members of the group identified by <group>.  Additionally
show only the named fields in the result.
`
}

func (p *ListMembersCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.ID, "group", "", "Name of the group to list")
	f.StringVar(&p.fields, "fields", "", "Fields to display")
}

func (p *ListMembersCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if p.ID == "" {
		fmt.Println("--group must be specified for group-members")
		return subcommands.ExitFailure
	}

	// Grab a client
	c, err := getClient()
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	// Obtain the membership list
	membersList, err := c.ListGroupMembers(p.ID)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	for _, m := range membersList {
		printEntity(m, p.fields)
	}

	return subcommands.ExitSuccess
}
