package ctl

import (
	"context"
	"flag"
	"fmt"
	"strings"

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

	entity, err := client.EntityInfo(serverAddr, serverPort, clientID, serviceID, p.ID)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	fields := []string{}
	if p.fields != "" {
		fields = strings.Split(p.fields, ",")
	} else {
		fields = []string{
			"ID",
			"uidNumber",
			"pgidNumber",
			"GECOS",
			"legalName",
			"displayName",
			"homedir",
			"shell",
			"graphicalShell",
		}
	}

	for _, f := range fields {
		switch f {
		case "ID":
			fmt.Printf("ID: %s\n", entity.GetID())
		case "uidNumber":
			fmt.Printf("uidNumber: %d\n", entity.GetUidNumber())
		case "pgidNumber":
			if entity.Meta != nil && entity.GetMeta().GetPgidNumber() != 0 {
				fmt.Printf("gidNumber: %d\n", entity.GetMeta().GetPgidNumber())
			}
		case "GECOS":
			if entity.Meta != nil && entity.GetMeta().GetGECOS() != "" {
				fmt.Printf("GECOS: %s\n", entity.GetMeta().GetGECOS())
			}
		case "legalName":
			if entity.Meta != nil && entity.GetMeta().GetLegalName() != "" {
				fmt.Printf("legalName: %s\n", entity.GetMeta().GetLegalName())
			}
		case "displayName":
			if entity.Meta != nil && entity.Meta.GetDisplayName() != "" {
				fmt.Printf("displayname: %s\n", entity.GetMeta().GetDisplayName())
			}
		case "homedir":
			if  entity.Meta != nil && entity.GetMeta().GetHomedir() != "" {
				fmt.Printf("homedir: %s\n", entity.GetMeta().GetHomedir())
			}
		case "shell":
			if entity.Meta != nil && entity.GetMeta().GetShell() != "" {
				fmt.Printf("shell: %s\n", entity.GetMeta().GetShell())
			}
		case "graphicalShell":
			if entity.Meta != nil && entity.GetMeta().GetGraphicalShell() != "" {
				fmt.Printf("graphicalShell: %s\n", entity.GetMeta().GetGraphicalShell())
			}
		}
	}
	return subcommands.ExitSuccess
}
