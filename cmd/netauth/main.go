package main

import (
	"context"
	"flag"
	"os"

	"github.com/NetAuth/NetAuth/internal/ctl"

	"github.com/google/subcommands"
)

var (
	serverAddr = flag.String("server", "localhost", "Server Address")
	serverPort = flag.Int("port", 8080, "Server port")
	clientID   = flag.String("client", "", "Client ID to send")
	serviceID  = flag.String("service", "netauthctl", "Service ID to send")
	entity     = flag.String("entity", "", "Entity to send in the request")
	secret     = flag.String("secret", "", "Secret to send in the request")
)

func main() {
	flag.Parse()

	// These are global options, they get passed inwards to
	// package level variables to the internal/ctl package so that
	// subcommands can access these without needing to redefine
	// them.  Any new variable here must get a package level
	// variable in internal/ctl/connopts.go and an associated
	// setter method in the same file.
	ctl.SetServerAddr(*serverAddr)
	ctl.SetServerPort(*serverPort)
	ctl.SetClientID(*clientID)
	ctl.SetServiceID(*serviceID)
	ctl.SetEntity(*entity)
	ctl.SetSecret(*secret)

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
