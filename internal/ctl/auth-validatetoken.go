package ctl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/NetAuth/NetAuth/pkg/client"
)

var (
	authValidateTokenCmd = &cobra.Command{
		Use:     "validate-token",
		Short:   "Validate a token with the server",
		Long:    authValidateTokenLongDocs,
		Example: authValidateTokenExample,
		Run:     authValidateTokenRun,
	}

	authValidateTokenLongDocs = `

validate-token sends the token to the server for validation.  The
server may perform additional scrutiny to satisfy the token's
legitimacy, and the result will be returned with the status of the
token.`

	authValidateTokenExample = `$ netauth auth validate-token
Token verified`
)

func init() {
	authCmd.AddCommand(authValidateTokenCmd)
}

func authValidateTokenRun(cmd *cobra.Command, args []string) {
	// Grab a client
	c, err := client.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Attempt to validate the token
	result, err := c.ValidateToken(viper.GetString("entity"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(result.GetMsg())
}
