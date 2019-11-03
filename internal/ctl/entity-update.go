package ctl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/NetAuth/NetAuth/pkg/netauth"

	pb "github.com/NetAuth/Protocol"
)

var (
	uEntity         string
	uPGroup         string
	uGECOS          string
	uLegalName      string
	uDisplayName    string
	uHomedir        string
	uShell          string
	uGraphicalShell string
	uBadgeNumber    string

	entityUpdateCmd = &cobra.Command{
		Use:     "update",
		Short:   "Update metadata on an entity",
		Long:    entityUpdateLongDocs,
		Example: entityUpdateExample,
		Args:    cobra.ExactArgs(1),
		Run:     entityUpdateRun,
	}

	entityUpdateLongDocs = `
The update command updates the typed metadata stored on an entity.
Fields are updated with the flags from this command, and are
overwritten with anything specified.
`

	entityUpdateExample = `netauth entity update demo2 --displayName "Demonstation User"
Metadata Updated
`
)

func init() {
	entityCmd.AddCommand(entityUpdateCmd)
	entityUpdateCmd.Flags().StringVar(&uPGroup, "primary-group", "", "Primary group")
	entityUpdateCmd.Flags().StringVar(&uGECOS, "GECOS", "", "GECOS")
	entityUpdateCmd.Flags().StringVar(&uLegalName, "legalName", "", "Legal name")
	entityUpdateCmd.Flags().StringVar(&uDisplayName, "displayName", "", "Display name")
	entityUpdateCmd.Flags().StringVar(&uHomedir, "homedir", "", "Home Directory")
	entityUpdateCmd.Flags().StringVar(&uShell, "shell", "", "User command interpreter")
	entityUpdateCmd.Flags().StringVar(&uGraphicalShell, "graphicalShell", "", "Graphical shell")
	entityUpdateCmd.Flags().StringVar(&uBadgeNumber, "badgeNumber", "", "Badge number")
}

func entityUpdateRun(cmd *cobra.Command, args []string) {
	uEntity = args[0]

	meta := &pb.EntityMeta{}
	if cmd.Flags().Changed("primary-group") {
		meta.PrimaryGroup = &uPGroup
	}
	if cmd.Flags().Changed("GECOS") {
		meta.GECOS = &uGECOS
	}
	if cmd.Flags().Changed("legalName") {
		meta.LegalName = &uLegalName
	}
	if cmd.Flags().Changed("displayName") {
		meta.DisplayName = &uDisplayName
	}
	if cmd.Flags().Changed("homedir") {
		meta.Home = &uHomedir
	}
	if cmd.Flags().Changed("shell") {
		meta.Shell = &uShell
	}
	if cmd.Flags().Changed("graphicalShell") {
		meta.GraphicalShell = &uGraphicalShell
	}
	if cmd.Flags().Changed("badgeNumber") {
		meta.BadgeNumber = &uBadgeNumber
	}

	ctx = netauth.Authorize(ctx, token())
	if err := rpc.EntityUpdate(ctx, uEntity, meta); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Metadata Updated")
}
