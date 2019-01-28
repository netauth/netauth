package ctl

import (
	"github.com/spf13/cobra"
)

var (
	cliCmd = &cobra.Command{
		Use: "cli",
		Short: "Extra utilities for the CLI",
	}
)

func init() {
	systemCmd.AddCommand(cliCmd)
}
