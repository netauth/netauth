package ctl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/netauth/netauth/pkg/netauth"
)

var (
	newGroupName        string
	newGroupNumber      int
	newGroupDisplayName string
	newGroupManagedBy   string

	groupCreateCmd = &cobra.Command{
		Use:     "create <name>",
		Short:   "Create a new group",
		Long:    groupCreateLongDocs,
		Example: groupCreateExample,
		Args:    cobra.ExactArgs(1),
		Run:     groupCreateRun,
	}

	groupCreateLongDocs = `
Create an group with the specified name.  Though there are no strict
requirements on the name beyond it being a single word that is
globally unique, it is strongly encouraged to make it exclusively of
lower case letters and numbers.  For the best compatibility, it is
recommended to start with a letter only.

Additional fields can be specified on the command line such as the
display name, or a group to defer management capability to.  If
desired a custom number can be provided, but the default behavior is
sufficient to select a valid unallocated number for the new group.

The caller must possess the CREATE_GROUP capability or be a GLOBAL_ROOT
operator for this command to succeed.`

	groupCreateExample = `$ netauth group create demo-group
New group created successfully`
)

func init() {
	groupCmd.AddCommand(groupCreateCmd)
	groupCreateCmd.Flags().IntVar(&newGroupNumber, "number", -1, "Number to assign.")
	groupCreateCmd.Flags().StringVar(&newGroupDisplayName, "display-name", "", "Group display name")
	groupCreateCmd.Flags().StringVar(&newGroupManagedBy, "managed-by", "", "Delegate management to this group")
}

func groupCreateRun(cmd *cobra.Command, args []string) {
	newGroupName = args[0]

	ctx = netauth.Authorize(ctx, token())

	if err := rpc.GroupCreate(ctx, newGroupName, newGroupDisplayName, newGroupManagedBy, newGroupNumber); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Group Created")
}
