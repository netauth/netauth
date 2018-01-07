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

func (p *PingCmd) SetFlags(f *flag.FlagSet) {}

func (p *PingCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	msg, err := client.Ping(serverAddr, serverPort, clientID)
	if err != nil {
		return subcommands.ExitFailure
	}

	fmt.Println(msg)
	return subcommands.ExitSuccess
}
