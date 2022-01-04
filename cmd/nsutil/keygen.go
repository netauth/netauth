package main

import (
	"github.com/spf13/cobra"
)

var (
	keygenCmd = &cobra.Command{
		Use:   "keygen",
		Short: "Generate keys of various types",
	}
)

func init() {
	rootCmd.AddCommand(keygenCmd)
}
