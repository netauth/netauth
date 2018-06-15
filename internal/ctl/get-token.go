package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
)

type GetTokenCmd struct{}

func (*GetTokenCmd) Name() string     { return "get-token" }
func (*GetTokenCmd) Synopsis() string { return "Obtain a token from a NetAuth server." }
func (*GetTokenCmd) Usage() string {
	return `get-token
  Attempt to obtain a token from the specified server.
`
}

func (*GetTokenCmd) SetFlags(f *flag.FlagSet) {}

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
