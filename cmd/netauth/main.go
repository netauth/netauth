package main

import (
	"context"
	"flag"
	"os"

	"github.com/NetAuth/NetAuth/internal/ctl"

	"github.com/google/subcommands"
)

func main() {
	flag.Parse()

	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&ctl.PingCmd{}, "")
	subcommands.Register(&ctl.AuthCmd{}, "")

	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
