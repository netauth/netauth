package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
)

// ValidateTokenCmd requests server side token validation.
type ValidateTokenCmd struct{}

// Name of this cmdlet is 'validate-token'
func (*ValidateTokenCmd) Name() string { return "validate-token" }

// Synopsis returns short-form usage information.
func (*ValidateTokenCmd) Synopsis() string { return "Validate an existing token with a NetAuth server." }

// Usage returns long-form usage information.
func (*ValidateTokenCmd) Usage() string {
	return `validate-token
  Send the token to the NetAuth server for validation.
`
}

// SetFlags is required by the interface but ValidateTokenCmd has no flags of its own.
func (*ValidateTokenCmd) SetFlags(f *flag.FlagSet) {}

// Execute runs this cmdlet.
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
