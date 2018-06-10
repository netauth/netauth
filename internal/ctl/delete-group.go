package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/NetAuth/NetAuth/pkg/client"

	"github.com/google/subcommands"
)

type DeleteGroupCmd struct {
	name        string
	displayName string
	gid         int
}

func (*DeleteGroupCmd) Name() string     { return "delete-group" }
func (*DeleteGroupCmd) Synopsis() string { return "Delete a group existing on the server." }
func (*DeleteGroupCmd) Usage() string {
	return `new-group --name <name>
Delete the named group.
`
}

func (p *DeleteGroupCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.name, "name", "", "Name for the new group.")
}

func (p *DeleteGroupCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
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

	result, err := c.DeleteGroup(p.name, t)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	fmt.Println(result.GetMsg())
	return subcommands.ExitSuccess
}
