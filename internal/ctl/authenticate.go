package ctl

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/NetAuth/NetAuth/pkg/client"

	"github.com/google/subcommands"
)

type AuthCmd struct {
	serverAddr string
	serverPort int
	clientID   string
	serviceID  string
	entity     string
	secret     string
}

func (*AuthCmd) Name() string     { return "auth" }
func (*AuthCmd) Synopsis() string { return "Authenticate to a NetAuth server." }
func (*AuthCmd) Usage() string {
	return `auth --entity <entity> --entity_secret <secret>
  Attempt to authenticate to the server specified.
`
}

func (p *AuthCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.serverAddr, "server", "localhost", "Server to connect to")
	f.IntVar(&p.serverPort, "port", 8080, "Server port")
	f.StringVar(&p.clientID, "client", "", "Client ID to send")
	f.StringVar(&p.serviceID, "service", "", "Service ID to send")
	f.StringVar(&p.entity, "entity", "", "Entity to authenticate as")
	f.StringVar(&p.secret, "secret", "", "Secret to authenticate with")
}

func (p *AuthCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if p.secret == "" {
		fmt.Print("Secret: ")
		_, err := fmt.Scanln(p.secret)
		if err != nil {
			log.Printf("Error: %s", err)
		}
	}
	ok := client.Authenticate(p.serverAddr, p.serverPort, p.clientID, p.serviceID, p.entity, p.secret)
	if !ok {
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}
