package ctl

import (
	"github.com/spf13/cobra"
)

var (
	authCmd = &cobra.Command{
		Use:   "auth <command>",
		Short: "Use and set authentication data",
		Long:  authCmdLongDocs,
	}

	authCmdLongDocs = `
The auth subystem deals in authentication information.  This includes
secrets and tokens.  Here you'll find the commands to change secrets,
chec them, get and destroy tokens and other intesting tasks around
authentication.`
)

func init() {
	rootCmd.AddCommand(authCmd)
}
