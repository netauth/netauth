package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
)

// EntityInfoCmd summons information on a named entity
type EntityInfoCmd struct {
	entityID string
	fields   string
}

// Name of this cmdlet is 'entity-info'
func (*EntityInfoCmd) Name() string { return "entity-info" }

// Synopsis for the cmdlet.
func (*EntityInfoCmd) Synopsis() string { return "Obtain information on an entity" }

// Usage info for the cmdlet.
func (*EntityInfoCmd) Usage() string {
	return `entity-info --entity <ID>  [--fields field1,field2...]
Print information about an entity.  The listed fields can be used to
limit the information that is printed.
`
}

// SetFlags processes the flags for this cmdlet.
func (p *EntityInfoCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.entityID, "entity", getEntity(), "ID to summon info for")
	f.StringVar(&p.fields, "fields", "", "Comma separated list of fields to display")
}

// Execute is called to run this cmdlet.
func (p *EntityInfoCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Grab a client
	c, err := getClient()
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	// Obtain entity info
	entity, err := c.EntityInfo(p.entityID)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	// Print the fields
	printEntity(entity, p.fields)

	return subcommands.ExitSuccess
}
