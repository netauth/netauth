package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
)

// CapabilitiesCmd supports the ModifyCapabilities RPC.
type CapabilitiesCmd struct {
	mode       string
	entityID   string
	groupName  string
	capability string
}

// Name of this cmdlet is 'modify-capabilities'
func (*CapabilitiesCmd) Name() string { return "modify-capabilities" }

// Synopsis for the cmdlet
func (*CapabilitiesCmd) Synopsis() string { return "Modify capabilities on an entity or group" }

// Usage of this cmdlet in long form.
func (*CapabilitiesCmd) Usage() string {
	return `modify-capabilities --capability <capability> <[--entity <ID>]|[--group <name>]> --mode <ADD|REMOVE>

Add or remove a capability from the named group or entity.  If both
are specififed (unsupported) then the group will be ignored.
`
}

// SetFlags is called to set flags specific to this cmdlet
func (p *CapabilitiesCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.mode, "mode", "ADD", "Mode, must be one of ADD or REMOVE")
	f.StringVar(&p.entityID, "entity", "", "Entity to modify")
	f.StringVar(&p.groupName, "group", "", "Group to modify")
	f.StringVar(&p.capability, "capability", "", "Capability to modify")
}

// Execute is the interface method that runs the actions of the cmdlet.
func (p *CapabilitiesCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if p.entityID == "" && p.groupName == "" {
		fmt.Println("Either --entity or --group must be specified!")
		return subcommands.ExitFailure
	}

	if p.mode != "ADD" && p.mode != "REMOVE" {
		fmt.Println("Mode must be either ADD or REMOVE")
		return subcommands.ExitFailure
	}

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

	result, err := c.ManageCapabilities(t, p.entityID, p.groupName, p.capability, p.mode)
	if result.GetMsg() != "" {
		fmt.Println(result.GetMsg())
	}
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}
