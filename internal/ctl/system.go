package ctl

import (
	"github.com/spf13/cobra"
)

var (
	systemCmd = &cobra.Command{
		Use:   "system",
		Short: "Internal system functions",
	}
)

func init() {
	rootCmd.AddCommand(systemCmd)
}
