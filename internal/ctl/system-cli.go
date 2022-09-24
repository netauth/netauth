package ctl

import (
	"github.com/spf13/cobra"
)

var (
	cliCmd = &cobra.Command{
		Use:   "cli",
		Short: "Extra utilities for the CLI",
		PersistentPreRun: cli_initialize,
	}
)

func init() {
	systemCmd.AddCommand(cliCmd)
}

func cli_initialize(*cobra.Command, []string) {}
