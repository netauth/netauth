package ctl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	authDestroyTokenCmd = &cobra.Command{
		Use:     "destroy-token",
		Short:   "Destroy an existing local token",
		Long:    authDestroyTokenLongDocs,
		Example: authDestroyTokenExample,
		Run:     authDestroyTokenRun,
	}

	authDestroyTokenLongDocs = `
destroy-token makes a best effort to remove the local token from the
system.  When this command returns the local token will either have
been destroyed or an error will be printed.  If this command returns
an error you cannot assume that the token has been removed!`

	authDestroyTokenExample = `$ netauth auth destroy-token
Token destroyed.`
)

func init() {
	authCmd.AddCommand(authDestroyTokenCmd)
}

func authDestroyTokenRun(cmd *cobra.Command, args []string) {
	// Destroy the token
	if err := rpc.DelToken(viper.GetString("entity")); err != nil {
		fmt.Printf("Error during token destruction: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Token destroyed.")
}
