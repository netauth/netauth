package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/NetAuth/NetAuth/pkg/client"

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

func (p *AuthCmd) SetFlags(f *flag.FlagSet) {}

func (p *AuthCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Ensure the secret has been obtained if it wasn't specified
	// on the command line.
	ensureSecret()

	// Grab a client
	c, err := client.New(serverAddr, serverPort, serviceID, clientID)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	// Attempt authentication
	msg, err := c.Authenticate(entity, secret)
	if err != nil {
		return subcommands.ExitFailure
	}
	fmt.Println(msg)
	return subcommands.ExitSuccess
}
