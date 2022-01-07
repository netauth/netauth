package ctl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	authInspectTokenCmd = &cobra.Command{
		Use:     "inspect-token",
		Short:   "Inspect a token locally",
		Long:    authInspectTokenLongDocs,
		Example: authInspectTokenExample,
		Run:     authInspectTokenRun,
	}

	authInspectTokenLongDocs = `
inspect-token prints a token for inspection locally.  Specifically it
prints the claims held in an encoded token.  Tokens are summoned on
demand, and this command will trigger an implicit call to get-token if
no local token is valid or available.  `

	authInspectTokenExample = `$ netauth auth inspect-token
Secret:
{root [GLOBAL_ROOT] 5}

$ netauth auth inspect-token
{root [GLOBAL_ROOT] 5}`
)

func init() {
	authCmd.AddCommand(authInspectTokenCmd)
}

func authInspectTokenRun(cmd *cobra.Command, args []string) {
	claims, err := tsvc.Validate(token())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("This token was issued to '%s'\n", claims.EntityID)
	if len(claims.EntityID) > 0 {
		fmt.Printf(" Capabilities:\n")
	} else {
		fmt.Printf("\n")
	}
	for i := range claims.Capabilities {
		fmt.Printf("  - %s\n", claims.Capabilities[i])
	}
}
