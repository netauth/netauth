package ctl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	systemPingCmd = &cobra.Command{
		Use:     "ping",
		Short:   "Ping the server and print the reply",
		Long:    systemPingLongDocs,
		Example: systemPingExample,
		Run:     systemPingRun,
	}

	systemPingLongDocs = `
The ping command provides an easy way to interogate a server and find
if it is behaving as expected.  The ping command requests a server to
pong back if with its health status.
`

	systemPingExample = `$ netauth system ping
NetAuth server on theGibson is ready to serve!`
)

func init() {
	systemCmd.AddCommand(systemPingCmd)
}

func systemPingRun(cmd *cobra.Command, args []string) {
	if err := rpc.SystemPing(ctx); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	fmt.Println("Pong")
}
