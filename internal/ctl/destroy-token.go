package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/NetAuth/NetAuth/pkg/client"

	"github.com/google/subcommands"
)

type DestroyTokenCmd struct{}

func (*DestroyTokenCmd) Name() string     { return "destroy-token" }
func (*DestroyTokenCmd) Synopsis() string { return "Destroy an existing local token." }
func (*DestroyTokenCmd) Usage() string {
	return `destroy-token
  Attempt to destroy the local authority token.  This command will
  make a best effort attempt to remove the local token.
`
}

func (*DestroyTokenCmd) SetFlags(f *flag.FlagSet) {}

func (*DestroyTokenCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Grab a client
	c, err := client.New(serverAddr, serverPort, serviceID, clientID)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	// Destroy the token
	if err := c.DestroyToken(entity); err != nil {
		fmt.Printf("Error during token destruction: %s\n", err)
		return subcommands.ExitFailure
	}

	fmt.Println("Token destroyed.")
	return subcommands.ExitSuccess
}
