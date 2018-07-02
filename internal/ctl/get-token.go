package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
)

// GetTokenCmd gets a NetAuth token for later use.
type GetTokenCmd struct{}

// Name of this cmdlet will be 'get-token'
func (*GetTokenCmd) Name() string { return "get-token" }

// Synopsis returns the short-form usage information
func (*GetTokenCmd) Synopsis() string { return "Obtain a token from a NetAuth server." }

// Usage returns the long-form usage information.
func (*GetTokenCmd) Usage() string {
	return `get-token
  Attempt to obtain a token from the specified server.
`
}

// SetFlags is required by the interface specification but GetTokenCmd
// takes no flags.
func (*GetTokenCmd) SetFlags(f *flag.FlagSet) {}

// Execute summons the token and stores it for later use.
func (*GetTokenCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Grab a client
	c, err := getClient()
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	// Attempt to get a token
	_, err = c.GetToken(getEntity(), getSecret())
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}
	fmt.Println("Token obtained")
	return subcommands.ExitSuccess
}
