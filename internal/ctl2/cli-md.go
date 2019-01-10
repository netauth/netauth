package ctl2

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var (
	cliMdCmd = &cobra.Command{
		Use:   "md <path>",
		Short: "Generate md pages at <path>",
		Args:  cobra.ExactArgs(1),
		Run:   cliMdRun,
	}
)

func init() {
	cliCmd.AddCommand(cliMdCmd)
}

func cliMdRun(cmd *cobra.Command, args []string) {
	err := doc.GenMarkdownTree(cmd.Root(), args[0])
	if err != nil {

		fmt.Println(err)
		os.Exit(1)
	}
}
