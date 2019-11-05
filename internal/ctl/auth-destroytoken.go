package ctl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/netauth/netauth/pkg/client"
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
	// Grab a client
	c, err := client.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Destroy the token
	if err := c.DestroyToken(viper.GetString("entity")); err != nil {
		fmt.Printf("Error during token destruction: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Token destroyed.")
}
