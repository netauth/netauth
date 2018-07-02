package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
)

// AuthCmd supports the AuthEntity RPC.
type AuthCmd struct{}

// Name of this cmdlet is 'auth'
func (*AuthCmd) Name() string { return "auth" }

// Synopsis for the cmdlet
func (*AuthCmd) Synopsis() string { return "Authenticate to a NetAuth server." }

// Usage of this cmdlet in long form.
func (*AuthCmd) Usage() string {
	return `auth
  Attempt to authenticate to the server specified.
`
}

// SetFlags is required by the interface but AuthCmd has no flags to
// set.
func (*AuthCmd) SetFlags(f *flag.FlagSet) {}

// Execute is the interface method that runs the actions of the cmdlet.
func (*AuthCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Grab a client
	c, err := getClient()
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	// Attempt authentication
	result, err := c.Authenticate(getEntity(), getSecret())
	if err != nil {
		return subcommands.ExitFailure
	}
	fmt.Println(result.GetMsg())
	return subcommands.ExitSuccess
}
