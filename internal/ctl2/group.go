package ctl2

import (
	"github.com/spf13/cobra"
)

var (
	groupCmd = &cobra.Command{
		Use:   "group",
		Short: "Manage groups and associated data",
	}
)

func init() {
	rootCmd.AddCommand(groupCmd)
}
