package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/NetAuth/NetAuth/internal/ctl"

	"github.com/google/subcommands"
)

func main() {
	flag.Parse()

	// Turn off the logging since the client should not be
	// spitting out any content unless its explicitly printed out
	// from internal/ctl.
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)

	// Register all the subcommands, each subcommand must be
	// registered after the builtins to be resolved in the right
	// order.  The order they are resolved here will not be the
	// order they are shown in the help output.
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&ctl.PingCmd{}, "System")
	subcommands.Register(&ctl.AuthCmd{}, "Authentication")
	subcommands.Register(&ctl.GetTokenCmd{}, "Authentication")
	subcommands.Register(&ctl.DestroyTokenCmd{}, "Authentication")
	subcommands.Register(&ctl.ValidateTokenCmd{}, "Authentication")
	subcommands.Register(&ctl.InspectTokenCmd{}, "Authentication")
	subcommands.Register(&ctl.ChangeSecretCmd{}, "Authentication")

	subcommands.Register(&ctl.NewEntityCmd{}, "Entity Administration")
	subcommands.Register(&ctl.RemoveEntityCmd{}, "Entity Administration")
	subcommands.Register(&ctl.EntityInfoCmd{}, "Entity Administration")
	subcommands.Register(&ctl.ModifyMetaCmd{}, "Entity Administration")
	subcommands.Register(&ctl.ModifyKeysCmd{}, "Entity Administration")

	subcommands.Register(&ctl.NewGroupCmd{}, "Group Administration")
	subcommands.Register(&ctl.DeleteGroupCmd{}, "Group Administration")
	subcommands.Register(&ctl.ListGroupsCmd{}, "Group Administration")
	subcommands.Register(&ctl.ModifyGroupCmd{}, "Group Administration")
	subcommands.Register(&ctl.GroupInfoCmd{}, "Group Administration")

	subcommands.Register(&ctl.ListMembersCmd{}, "Membership Administration")
	subcommands.Register(&ctl.EntityMembershipCmd{}, "Membership Administration")
	subcommands.Register(&ctl.GroupExpansionsCmd{}, "Membership Administration")

	subcommands.Register(&ctl.CapabilitiesCmd{}, "Capabilities Administration")

	// Register in the global flags as important
	subcommands.ImportantFlag("server")
	subcommands.ImportantFlag("port")
	subcommands.ImportantFlag("client")
	subcommands.ImportantFlag("service")
	subcommands.ImportantFlag("entity")
	subcommands.ImportantFlag("secret")

	// By default we will run the functions at background context.
	// Below  this call  level it  may be  necessary to  reset the
	// context, but the initial call level can be background.
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
