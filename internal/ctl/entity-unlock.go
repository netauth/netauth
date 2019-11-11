package ctl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/netauth/netauth/pkg/netauth"
)

var (
	entityUnlockCmd = &cobra.Command{
		Use:     "unlock <ID>",
		Short:   "Unlock the entity with the specified ID",
		Long:    entityUnlockLongDocs,
		Example: entityUnlockExample,
		Args:    cobra.ExactArgs(1),
		Run:     entityUnlockRun,
	}

	entityUnlockLongDocs = `
Unlock an entity with the specified ID.  A locked entity cannot
authenticate successfully, even when presenting the correct secret.

The caller must possess the UNLOCK_ENTITY capability or be a
GLOBAL_ROOT operator for this command to succeed.`

	entityUnlockExample = `$ netauth entity lock demo
Entity is now unlocked
`
)

func init() {
	entityCmd.AddCommand(entityUnlockCmd)
}

func entityUnlockRun(cmd *cobra.Command, args []string) {
	ctx = netauth.Authorize(ctx, token())

	if err := rpc.EntityUnlock(ctx, args[0]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Entity Unlocked")
}
