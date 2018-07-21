package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
)

// NewEntityCmd requests entity creation on the server.
type NewEntityCmd struct {
	ID     string
	number int
	secret string
}

// Name of this cmdlet is 'new-entity'
func (*NewEntityCmd) Name() string { return "new-entity" }

// Synopsis returns the short-form usage information.
func (*NewEntityCmd) Synopsis() string { return "Add a new entity to the server" }

// Usage returns the long-form usage information.
func (*NewEntityCmd) Usage() string {
	return `new-entity --ID <ID> --number <number> --secret <secret>
  Create a new entity with the specified ID, number, and secret.
  number may be ommitted to select the next available number.
  Secret may be ommitted to leave unset.`
}

// SetFlags sets the flags specific to this command.
func (p *NewEntityCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.ID, "ID", "", "ID for the new entity")
	f.IntVar(&p.number, "number", -1, "number for the new entity")
	f.StringVar(&p.secret, "secret", "", "secret for the new entity")
}

// Execute runs the cmdlet.
func (p *NewEntityCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
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
	result, err := c.NewEntity(p.ID, number, p.secret, t)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}
	fmt.Println(result.GetMsg())

	return subcommands.ExitSuccess
}
