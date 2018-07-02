package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
)

// PingCmd requests the server to run its health checks and return the status.
type PingCmd struct{}

// Name of this cmdlet is 'ping'
func (*PingCmd) Name() string { return "ping" }

// Synopsis returns the short-form usage information.
func (*PingCmd) Synopsis() string { return "Ping a NetAuth server." }

// Usage returns the long-form usage inforamtion.
func (*PingCmd) Usage() string {
	return `ping:
  Ping the server and print the reply.
`
}

// SetFlags is required by the interface, but PingCmd has no flags of its own.
func (*PingCmd) SetFlags(f *flag.FlagSet) {}

// Execute runs the cmdlet.
func (*PingCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Grab a client
	c, err := getClient()
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	result, err := c.Ping()
	if err != nil {
		return subcommands.ExitFailure
	}

	fmt.Println(result.GetMsg())
	if !result.GetHealthy() {
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}
