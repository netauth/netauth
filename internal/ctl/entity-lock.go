package ctl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/netauth/netauth/pkg/netauth"
)

var (
	entityLockCmd = &cobra.Command{
		Use:   "lock <ID>",
		Short: "Lock the entity with the specified ID",
		Long:  entityLockLongDocs,
		Args:  cobra.ExactArgs(1),
		Run:   entityLockRun,
	}

	entityLockLongDocs = `
Lock an entity with the specified ID.  A locked entity cannot
authenticate successfully, even when presenting the correct secret.

The caller must posess the LOCK_ENTITY capability or be a GLOBAL_ROOT
operator for this command to succeed.`

	entityLockExample = `$ netauth entity lock demo
Entity is now locked
`
)

func init() {
	entityCmd.AddCommand(entityLockCmd)
}

func entityLockRun(cmd *cobra.Command, args []string) {
	ctx = netauth.Authorize(ctx, token())

	if err := rpc.EntityLock(ctx, args[0]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Entity Locked")
}
