package ctl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cliBashCmd = &cobra.Command{
		Use:   "bash <path>",
		Short: "Generate bash completions at <path>",
		Args:  cobra.ExactArgs(1),
		Run:   cliBashRun,
	}
)

func init() {
	cliCmd.AddCommand(cliBashCmd)
}

func cliBashRun(cmd *cobra.Command, args []string) {
	err := cmd.Root().GenBashCompletionFile(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
