package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/NetAuth/NetAuth/pkg/client"

	"github.com/google/subcommands"
)

type NewEntityCmd struct {
	ID        string
	uidNumber int
	secret    string
}

func (*NewEntityCmd) Name() string     { return "new-entity" }
func (*NewEntityCmd) Synopsis() string { return "Add a new entity to the server" }
func (*NewEntityCmd) Usage() string {
	return `new-entity --ID <ID> --uidNumber <number> --secret <secret>
  Create a new entity with the specified ID, uidNumber, and secret.
  uidNumber may be ommitted to select the next available uidNumber.
  Secret may be ommitted to leave unset.`
}

func (p *NewEntityCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.ID, "ID", "", "ID for the new entity")
	f.IntVar(&p.uidNumber, "uidNumber", -1, "uidNumber for the new entity")
	f.StringVar(&p.secret, "secret", "", "secret for the new entity")
}

func (p *NewEntityCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
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

	// The uidNumber has to be an int32 to be accepted into the
	// system.  This is for reasons related to protobuf.
	uidNumber := int32(p.uidNumber)
	msg, err := c.NewEntity(p.ID, uidNumber, p.secret, t)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	fmt.Println(msg)
	return subcommands.ExitSuccess
}
