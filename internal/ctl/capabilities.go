package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"

	"github.com/NetAuth/NetAuth/pkg/client"
)

type CapabilitiesCmd struct {
	mode       string
	entity     string
	group      string
	capability string
}

func (*CapabilitiesCmd) Name() string     { return "modify-capabilities" }
func (*CapabilitiesCmd) Synopsis() string { return "Modify capabilities on an entity or group" }
func (*CapabilitiesCmd) Usage() string {
	return `modify-capabilities --capability <capability> <[--entity <ID>]|[--group <name>]> --mode <ADD|REMOVE>

Add or remove a capability from the named group or entity.  If both
are specififed (unsupported) then the group will be ignored.
`
}

func (p *CapabilitiesCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.mode, "mode", "ADD", "Mode, must be one of ADD or REMOVE")
	f.StringVar(&p.entity, "entity", "", "Entity to modify")
	f.StringVar(&p.group, "group", "", "Group to modify")
	f.StringVar(&p.capability, "capability", "", "Capability to modify")
}

func (p *CapabilitiesCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if p.entity == "" && p.group == "" {
		fmt.Println("Either --entity or --group must be specified!")
		return subcommands.ExitFailure
	}

	if p.mode != "ADD" && p.mode != "REMOVE" {
		fmt.Println("Mode must be either ADD or REMOVE")
		return subcommands.ExitFailure
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

	result, err := c.ManageCapabilities(t, p.entity, p.group, p.capability, p.mode)
	if result.GetMsg() != "" {
		fmt.Println(result.GetMsg())
	}
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}
