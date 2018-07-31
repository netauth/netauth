package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
)

// CreateEntityCmd requests entity creation on the server.
type CreateEntityCmd struct {
	entityID string
	number   int
	secret   string
}

// Name of this cmdlet is 'new-entity'
func (*CreateEntityCmd) Name() string { return "create-entity" }

// Synopsis returns the short-form usage information.
func (*CreateEntityCmd) Synopsis() string { return "Add a new entity to the server" }

// Usage returns the long-form usage information.
func (*CreateEntityCmd) Usage() string {
	return `create-entity --ID <ID> --number <number> --secret <secret>
  Create a new entity with the specified ID, number, and secret.
  number may be ommitted to select the next available number.
  Secret may be ommitted to leave unset.`
}

// SetFlags sets the flags specific to this command.
func (p *CreateEntityCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.entityID, "ID", "", "ID for the new entity")
	f.IntVar(&p.number, "number", -1, "number for the new entity")
	f.StringVar(&p.secret, "secret", "", "secret for the new entity")
}

// Execute runs the cmdlet.
func (p *CreateEntityCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
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

	// The number has to be an int32 to be accepted into the
	// system.  This is for reasons related to protobuf.
	number := int32(p.number)
	result, err := c.NewEntity(p.entityID, number, p.secret, t)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}
	fmt.Println(result.GetMsg())

	return subcommands.ExitSuccess
}
