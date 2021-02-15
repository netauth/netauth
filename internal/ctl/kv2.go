package ctl

import (
	"github.com/spf13/cobra"
)

var (
	kv2Cmd = &cobra.Command{
		Use:   "kv2",
		Short: "Manage KV2 Data",
	}
)

func init() {
	rootCmd.AddCommand(kv2Cmd)
}
