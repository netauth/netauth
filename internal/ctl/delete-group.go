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
	displayName    string
	gid int
}

func (*DeleteGroupCmd) Name() string     { return "delete-group" }
func (*DeleteGroupCmd) Synopsis() string { return "Delete a group existing on the server." }
func (*DeleteGroupCmd) Usage() string {
	return `new-group --name <name>
Delete the named group.
`}

func (p *DeleteGroupCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.name, "name", "", "Name for the new group.")
}

func (p *DeleteGroupCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Ensure that the secret has been obtained to authorize this
	// command
	ensureSecret()

	// Grab a client
	c, err := client.New(serverAddr, serverPort, serviceID, clientID)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	msg, err := c.DeleteGroup(entity, secret, p.name)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	fmt.Println(msg)
	return subcommands.ExitSuccess
}
