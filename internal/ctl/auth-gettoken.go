package ctl

import (
	"fmt"

	"github.com/spf13/cobra"
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
	// Attempt to get a token
	refreshToken()
	fmt.Println("Token obtained")
}
