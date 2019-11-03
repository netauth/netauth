package ctl

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var (
	systemStatusCmd = &cobra.Command{
		Use:     "status",
		Short:   "Request a status report from the server",
		Long:    systemStatusLongDocs,
		Example: systemStatusExample,
		Run:     systemStatusRun,
	}

	systemStatusLongDocs = `
The status command provides an easy way to interogate a server and
find if it is behaving as expected.  The status command requests a
server to return information on each subsystem and to report failure
information if any subsystem is reporting an unhealthy state.
`

	systemStatusExample = `$ netauth system status`
)

func init() {
	systemCmd.AddCommand(systemStatusCmd)
}

func systemStatusRun(cmd *cobra.Command, args []string) {
	res, err := rpc.SystemStatus(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	fmt.Println("Server Status:", boolToPassFail(res.GetSystemOK()))
	fmt.Println()

	if res.GetFirstFailure() != nil {
		fmt.Println("The first system to report fail status:")
		fmt.Println(res.GetFirstFailure())
		fmt.Println()
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, s := range res.GetSubSystems() {
		fmt.Fprintf(w, "%s\t%s\t%s\n", boolToPassFail(s.GetOK()), s.GetName(), s.GetFaultMessage())
	}
	w.Flush()
}

func boolToPassFail(b bool) string {
	if !b {
		return "[FAIL]"
	}
	return "[PASS]"
}
