package ctl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/NetAuth/NetAuth/pkg/client"

	pb "github.com/NetAuth/Protocol"
)

var (
	uGDisplayName string
	uGManagedBy   string

	groupUpdateCmd = &cobra.Command{
		Use:     "update",
		Short:   "Update metadata on an group",
		Long:    groupUpdateLongDocs,
		Example: groupUpdateExample,
		Args:    cobra.ExactArgs(1),
		Run:     groupUpdateRun,
	}

	groupUpdateLongDocs = `
The update command updates the typed metadata stored on an group.
Fields are updated with the flags from this command, and are
overwritten with anything specified.
`

	groupUpdateExample = `netuath group update example-group --display-name "Example Group"
Group modified successfully
`
)

func init() {
	groupCmd.AddCommand(groupUpdateCmd)
	groupUpdateCmd.Flags().StringVar(&uGDisplayName, "display-name", "", "Display Name")
	groupUpdateCmd.Flags().StringVar(&uGManagedBy, "managed-by", "", "Dlegated management group")
}

func groupUpdateRun(cmd *cobra.Command, args []string) {
	grp := &pb.Group{Name: &args[0]}
	if cmd.Flags().Changed("display-name") {
		grp.DisplayName = &uGDisplayName
	}
	if cmd.Flags().Changed("managed-by") {
		grp.ManagedBy = &uGManagedBy
	}

	// Grab a client
	c, err := client.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get the authorization token
	t, err := getToken(c, viper.GetString("entity"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	result, err := c.ModifyGroupMeta(grp, t)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(result.GetMsg())
}
