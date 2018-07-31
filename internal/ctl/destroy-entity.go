package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
)

// DestroyEntityCmd requests the server to remove an entity.
type DestroyEntityCmd struct {
	entityID string
}

// Name of this cmdlet is 'remove-entity'
func (*DestroyEntityCmd) Name() string { return "destroy-entity" }

// Synopsis returns the short-form usage information.
func (*DestroyEntityCmd) Synopsis() string { return "Remove an entity from the server" }

// Usage returns the long-form usage information.
func (*DestroyEntityCmd) Usage() string {
	return `destroy-entity --entity <ID>
Remove the specified entity from the server.`
}

// SetFlags sets the cmdlet specific flags.
func (p *DestroyEntityCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.entityID, "entity", "", "ID for the entity to be removed")
}

// Execute runs the cmdlet
func (p *DestroyEntityCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Grab a client
	c, err := getClient()
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	// Get the authorization token
	t, err := getToken(c, getEntity())
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	// Remove the entity
	result, err := c.RemoveEntity(p.entityID, t)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	fmt.Println(result.GetMsg())
	return subcommands.ExitSuccess
}
