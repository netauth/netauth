package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/NetAuth/NetAuth/pkg/client"

	"github.com/google/subcommands"
)

type GroupInfoCmd struct {
	name   string
	fields string
}

func (*GroupInfoCmd) Name() string     { return "group-info" }
func (*GroupInfoCmd) Synopsis() string { return "Obtain information on a group" }
func (*GroupInfoCmd) Usage() string {
	return `group-info --name <name> [--fields field1,field2,field3...]

Return the fields of a group.  This will provide information on a
single group, as opposed to attempting to list all groups.
`
}

func (p *GroupInfoCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.fields, "fields", "", "Comma seperated list of fields to display")
	f.StringVar(&p.name, "name", "", "Name of the group to query")
}

func (p *GroupInfoCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Grab a client
	c, err := client.New(serverAddr, serverPort, serviceID, clientID)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	// Obtain group info
	result, err := c.GroupInfo(p.name)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	printGroup(result.GetGroup(), p.fields)

	if len(result.GetManaged()) > 0 {
		fmt.Printf("The following group(s) are managed by %s\n", p.name)
	}
	for _, gn := range result.GetManaged() {
		fmt.Printf("  - %s\n", gn)
	}
	return subcommands.ExitSuccess
}
