package ctl

import (
	"context"
	"flag"
	"fmt"
	"log"

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
	if secret == "" {
		fmt.Print("Secret: ")
		_, err := fmt.Scanln(&secret)
		if err != nil {
			log.Printf("Error: %s", err)
		}
	}
	// Authenticate to the server, the variables that come from
	// "nowhere" are package-scoped and originate in connopts.go
	// adjacent to this file.
	ok := client.Authenticate(serverAddr, serverPort, clientID, serviceID, entity, secret)
	if !ok {
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}
