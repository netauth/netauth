package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/NetAuth/NetAuth/pkg/client"

	"github.com/google/subcommands"
)

type PingCmd struct{}

func (*PingCmd) Name() string     { return "ping" }
func (*PingCmd) Synopsis() string { return "Ping a NetAuth server." }
func (*PingCmd) Usage() string {
	return `ping:
  Ping the server and print the reply.
`
}

func (*PingCmd) SetFlags(f *flag.FlagSet) {}

func (*PingCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Grab a client
	c, err := client.New(serverAddr, serverPort, serviceID, clientID)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	result, err := c.Ping()
	if err != nil {
		return subcommands.ExitFailure
	}

	fmt.Println(result.GetMsg())
	if result.GetHealthy() {
		return subcommands.ExitSuccess
	} else {
		return subcommands.ExitFailure
	}
}
