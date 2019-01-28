package ctl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cliZshCmd = &cobra.Command{
		Use:   "zsh <path>",
		Short: "Generate zsh completions at <path>",
		Args:  cobra.ExactArgs(1),
		Run:   cliZshRun,
	}
)

func init() {
	cliCmd.AddCommand(cliZshCmd)
}

func cliZshRun(cmd *cobra.Command, args []string) {
	err := cmd.Root().GenZshCompletionFile(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
