package ctl

import (
	"context"
	"flag"

	"github.com/NetAuth/NetAuth/pkg/client"

	"github.com/google/subcommands"
)

type PingCmd struct {
	serverAddr string
	serverPort int
	clientID   string
}

func (*PingCmd) Name() string     { return "ping" }
func (*PingCmd) Synopsis() string { return "Ping a NetAuth server." }
func (*PingCmd) Usage() string {
	return `ping:
  Ping the server and print the reply.
`
}

func (p *PingCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.serverAddr, "server", "localhost", "Server to connect to")
	f.IntVar(&p.serverPort, "port", 8080, "Server port")
	f.StringVar(&p.clientID, "client", "", "Client ID to send")
}

func (p *PingCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	ok := client.Ping(p.serverAddr, p.serverPort, p.clientID)
	if !ok {
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}
