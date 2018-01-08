package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/NetAuth/NetAuth/pkg/client"

	"github.com/google/subcommands"
)

type RemoveEntityCmd struct {
	ID string
}

func (*RemoveEntityCmd) Name() string     { return "remove-entity" }
func (*RemoveEntityCmd) Synopsis() string { return "Add a remove entity to the server" }
func (*RemoveEntityCmd) Usage() string {
	return `remove-entity --ID <ID>
Remove the specified entity from the server.`
}

func (p *RemoveEntityCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.ID, "ID", "", "ID for the entity to be removed")
}

func (p *RemoveEntityCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Ensure that the secret has been obtained to authorize this
	// command
	ensureSecret()

	msg, err := client.RemoveEntity(serverAddr, serverPort, clientID, serviceID, entity, secret, p.ID)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	fmt.Println(msg)
	return subcommands.ExitSuccess
}
