package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/NetAuth/NetAuth/pkg/client"

	"github.com/google/subcommands"
)

type EntityInfoCmd struct {
	ID     string
	fields string
}

func (*EntityInfoCmd) Name() string     { return "entity-info" }
func (*EntityInfoCmd) Synopsis() string { return "Obtain information on an entity" }
func (*EntityInfoCmd) Usage() string {
	return `entity-info --ID <ID>  [--fields field1,field2...]
Print information about an entity.  The listed fields can be used to
limit the information that is printed`
}

func (p *EntityInfoCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.ID, "ID", "", "ID for the new entity")
	f.StringVar(&p.fields, "secret", "", "secret for the new entity")
}

func (p *EntityInfoCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
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

	// Obtain entity info
	entity, err := c.EntityInfo(p.ID)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	// Print the fields
	printEntity(entity, p.fields)

	return subcommands.ExitSuccess
}
