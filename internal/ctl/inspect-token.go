package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/NetAuth/NetAuth/pkg/client"

	"github.com/google/subcommands"
)

type InspectTokenCmd struct{}

func (*InspectTokenCmd) Name() string     { return "inspect-token" }
func (*InspectTokenCmd) Synopsis() string { return "Inspect an existing token locally." }
func (*InspectTokenCmd) Usage() string {
	return `validate-token
  Inspect the token locally, printing its contents if it is valid.
`
}

func (*InspectTokenCmd) SetFlags(f *flag.FlagSet) {}

func (*InspectTokenCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Grab a client
	c, err := client.New(serverAddr, serverPort, serviceID, clientID)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	// Attempt to validate the token
	t, err := c.GetToken(entity, secret)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	claims, err := c.InspectToken(t)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}
	fmt.Println(claims)
	return subcommands.ExitSuccess
}
