package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
)

// RemoveEntityCmd requests the server to remove an entity.
type RemoveEntityCmd struct {
	entityID string
}

// Name of this cmdlet is 'remove-entity'
func (*RemoveEntityCmd) Name() string { return "remove-entity" }

// Synopsis returns the short-form usage information.
func (*RemoveEntityCmd) Synopsis() string { return "Add a remove entity to the server" }

// Usage returns the long-form usage information.
func (*RemoveEntityCmd) Usage() string {
	return `remove-entity --entity <ID>
Remove the specified entity from the server.`
}

// SetFlags sets the cmdlet specific flags.
func (p *RemoveEntityCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.entityID, "entity", "", "ID for the entity to be removed")
}

// Execute runs the cmdlet
func (p *RemoveEntityCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
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
