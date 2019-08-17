package ctl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/NetAuth/NetAuth/pkg/client"
)

var (
	unlockEntityID string

	entityUnlockCmd = &cobra.Command{
		Use:   "unlock <ID>",
		Short: "Unlock the entity with the specified ID",
		Long:  entityUnlockLongDocs,
		Args:  cobra.ExactArgs(1),
		Run:   entityUnlockRun,
	}

	entityUnlockLongDocs = `
Unlock an entity with the specified ID.  A locked entity cannot
authenticate successfully, even when presenting the correct secret.

The caller must posess the UNLOCK_ENTITY capability or be a
GLOBAL_ROOT operator for this command to succeed.`

	entityUnlockExample = `$ netauth entity lock demo
Entity is now unlocked
`)

func init() {
	entityCmd.AddCommand(entityUnlockCmd)
}

func entityUnlockRun(cmd *cobra.Command, args []string) {
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

	result, err := c.UnlockEntity(t, lockEntityID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(result.GetMsg())
}
