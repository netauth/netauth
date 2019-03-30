package ctl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/NetAuth/NetAuth/pkg/client"
)

var (
	authGetTokenCmd = &cobra.Command{
		Use:   "get-token",
		Short: "Request a new token from the server",
		Long:  authGetTokenLongDocs,
		Run:   authGetTokenRun,
	}

	authGetTokenLongDocs = `
get-token retrieves a token from the server if one is not already
available locally.  If a token is available locally and is still
valid, the server will not be contacted.`

	authGetTokenExample = `$ netauth auth get-token
Secret:
Token obtained`
)

func init() {
	authCmd.AddCommand(authGetTokenCmd)
}

func authGetTokenRun(cmd *cobra.Command, args []string) {
	// Grab a client
	c, err := client.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Attempt to get a token
	_, err = getToken(c, viper.GetString("entity"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Token obtained")
}
