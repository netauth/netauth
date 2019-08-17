package ctl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/NetAuth/NetAuth/pkg/client"
)

var (
	lockEntityID string

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
`)

func init() {
	entityCmd.AddCommand(entityLockCmd)
}

func entityLockRun(cmd *cobra.Command, args []string) {
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

	lockEntityID = args[0]

	result, err := c.LockEntity(t, lockEntityID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(result.GetMsg())
}
