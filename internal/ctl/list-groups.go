package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/NetAuth/NetAuth/pkg/client"

	"github.com/google/subcommands"
)

type ListGroupsCmd struct {
	fields string
}

func (*ListGroupsCmd) Name() string     { return "list-groups" }
func (*ListGroupsCmd) Synopsis() string { return "Obtain a list of groups" }
func (*ListGroupsCmd) Usage() string {
	return `list-groups [--fields field1,field2,field3...]
This command will return a list of groups, additional formatting
options can be selected for additional information.
`
}

func (p *ListGroupsCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.fields, "fields", "", "Comma seperated list of fields to display")
}

func (p *ListGroupsCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Grab a client
	c, err := client.New(serverAddr, serverPort, serviceID, clientID)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	// Obtain group list
	list, err := c.ListGroups()
	if err != nil {
		return subcommands.ExitFailure
	}

	// Print the list
	for _, g := range list {
		printGroup(g, p.fields)
	}

	return subcommands.ExitSuccess
}
