package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
)

// DestroyTokenCmd clears the local token.
type DestroyTokenCmd struct{}

// Name to return for this cmdlet.
func (*DestroyTokenCmd) Name() string { return "destroy-token" }

// Synopsis for the cmdlet.
func (*DestroyTokenCmd) Synopsis() string { return "Destroy an existing local token." }

// Usage for the cmdlet.
func (*DestroyTokenCmd) Usage() string {
	return `destroy-token
  Attempt to destroy the local authority token.  This command will
  make a best effort attempt to remove the local token.
`
}

// SetFlags is required by the interface but DestroyTokenCmd has no
// flags to set.
func (*DestroyTokenCmd) SetFlags(f *flag.FlagSet) {}

// Execute is the interface method that runs the actions of the cmdlet.
func (*DestroyTokenCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Grab a client
	c, err := getClient()
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	// Destroy the token
	if err := c.DestroyToken(getEntity()); err != nil {
		fmt.Printf("Error during token destruction: %s\n", err)
		return subcommands.ExitFailure
	}

	fmt.Println("Token destroyed.")
	return subcommands.ExitSuccess
}
