package ctl2

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/NetAuth/NetAuth/pkg/client"
)

var (
	authCheckCmd = &cobra.Command{
		Use:     "check",
		Short:   "Check authentication credentials",
		Long:    authCheckLongDocs,
		Example: authCheckExample,
		Run:     authCheckRun,
	}

	authCheckLongDocs = `
The check command can be used to check authentication values without
requesting a token.  This command simply sends the values to the
server and returns the status from the server with no other
processing.  The entity that is checked can be influenced with the
global entity flag.`

	authCheckExample = `$ netauth auth check
Secret:
Entity authentication succeeded`
)

func init() {
	authCmd.AddCommand(authCheckCmd)
}

func authCheckRun(cmd *cobra.Command, args []string) {
	// Grab a client
	c, err := client.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Attempt authentication
	result, err := c.Authenticate(viper.GetString("entity"), getSecret(""))
	if err != nil {
		os.Exit(1)
	}
	fmt.Println(result.GetMsg())
}
