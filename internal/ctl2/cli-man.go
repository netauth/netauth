package ctl2

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var (
	cliManCmd = &cobra.Command{
		Use:   "man <path>",
		Short: "Generate man pages at <path>",
		Args:  cobra.ExactArgs(1),
		Run:   cliManRun,
	}
)

func init() {
	cliCmd.AddCommand(cliManCmd)
}

func cliManRun(cmd *cobra.Command, args []string) {
	header := &doc.GenManHeader{
		Title:   "NETAUTH",
		Section: "3",
	}
	err := doc.GenManTree(cmd.Root(), header, args[0])
	if err != nil {

		fmt.Println(err)
		os.Exit(1)
	}
}
