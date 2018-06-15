package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
)

type AuthCmd struct{}

func (*AuthCmd) Name() string     { return "auth" }
func (*AuthCmd) Synopsis() string { return "Authenticate to a NetAuth server." }
func (*AuthCmd) Usage() string {
	return `auth
  Attempt to authenticate to the server specified.
`
}

func (*AuthCmd) SetFlags(f *flag.FlagSet) {}

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
