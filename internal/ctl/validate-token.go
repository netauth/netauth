package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
)

type ValidateTokenCmd struct{}

func (*ValidateTokenCmd) Name() string     { return "validate-token" }
func (*ValidateTokenCmd) Synopsis() string { return "Validate an existing token with a NetAuth server." }
func (*ValidateTokenCmd) Usage() string {
	return `validate-token
  Send the token to the NetAuth server for validation.
`
}

func (*ValidateTokenCmd) SetFlags(f *flag.FlagSet) {}

func (*ValidateTokenCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Grab a client
	c, err := getClient()
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	// Attempt to validate the token
	result, err := c.ValidateToken(getEntity())
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}
	fmt.Println(result.GetMsg())
	return subcommands.ExitSuccess
}
