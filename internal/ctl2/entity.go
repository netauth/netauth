package ctl2

import (
	"github.com/spf13/cobra"
)

var (
	entityCmd = &cobra.Command{
		Use:   "entity",
		Short: "Manage entities and associated data",
	}
)

func init() {
	rootCmd.AddCommand(entityCmd)
}
