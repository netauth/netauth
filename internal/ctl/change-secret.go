package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/NetAuth/NetAuth/pkg/client"

	"github.com/google/subcommands"
)

type ChangeSecretCmd struct {
	ID     string
	secret string
}

func (*ChangeSecretCmd) Name() string     { return "change-secret" }
func (*ChangeSecretCmd) Synopsis() string { return "Change the secret for a given entity" }
func (*ChangeSecretCmd) Usage() string {
	return `change-secret --ID <ID>  --secret <secret>
Change the secret for the listed entity.  If no entity is provided the
entity specified by the top level flag will be used instead.`
}

func (p *ChangeSecretCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.ID, "ID", "", "ID for the new entity")
	f.StringVar(&p.secret, "secret", "", "secret for the new entity")
}

func (p *ChangeSecretCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// If the entity wasn't provided, use the one that was set
	// earlier.
	if p.ID == "" {
		p.ID = entity
	}

	// Grab a client
	c, err := client.New(serverAddr, serverPort, serviceID, clientID)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	// Get the authorization token
	t, err := c.GetToken(entity, secret)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	// Change the secret
	msg, err := c.ChangeSecret(entity, secret, p.ID, p.secret, t)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	fmt.Println(msg)
	return subcommands.ExitSuccess
}
