package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
)

// InspectTokenCmd examines the local token and prints properties about it
type InspectTokenCmd struct{}

// Name of this cmdlet is 'inspect-token'
func (*InspectTokenCmd) Name() string { return "inspect-token" }

// Synopsis returns the short-form usage.
func (*InspectTokenCmd) Synopsis() string { return "Inspect an existing token locally." }

// Usage returns the long-form usage.
func (*InspectTokenCmd) Usage() string {
	return `inspect-token
  Inspect the token locally, printing its contents if it is valid.
`
}

// SetFlags is required by the interface, but InspectTokenCmd has no flags of its own.
func (*InspectTokenCmd) SetFlags(f *flag.FlagSet) {}

// Execute runs the cmdlet.
func (*InspectTokenCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Grab a client
	c, err := getClient()
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	// Attempt to validate the token
	t, err := c.GetToken(getEntity(), getSecret())
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
