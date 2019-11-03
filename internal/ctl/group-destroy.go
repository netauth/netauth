package ctl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/NetAuth/NetAuth/pkg/netauth"
)

var (
	groupDestroyCmd = &cobra.Command{
		Use:     "destroy <name>",
		Short:   "Destroy an existing group",
		Long:    groupDestroyLongDocs,
		Example: groupDestroyExample,
		Args:    cobra.ExactArgs(1),
		Run:     groupDestroyRun,
	}

	groupDestroyLongDocs = `
Destroy the group with the specified name.  The group is deleted
immediately and without confirmation, please ensure you have typed the
ID correctly.

Referential integrity is not checked before deletion.  You are
strongly encouraged to empty groups before deleting them as well as
remove any expansions that target the group to be deleted.

The caller must posess the DESTROY_GROUP capability or be a
GLOBAL_ROOT operator for this command to succeed.
`

	groupDestroyExample = `$ netauth group destroy demo-group
Group removed successfully`
)

func init() {
	groupCmd.AddCommand(groupDestroyCmd)
}

func groupDestroyRun(cmd *cobra.Command, args []string) {
	ctx = netauth.Authorize(ctx, token())

	if err := rpc.GroupDestroy(ctx, args[0]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Group Destroyed")
}
