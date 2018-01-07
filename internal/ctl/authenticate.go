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
	// Authenticate to the server, the variables that come from
	// "nowhere" are package-scoped and originate in connopts.go
	// adjacent to this file.
	msg, err := client.Authenticate(serverAddr, serverPort, clientID, serviceID, entity, secret)
	if err != nil {
		return subcommands.ExitFailure
	}
	fmt.Println(msg)
	return subcommands.ExitSuccess
}
