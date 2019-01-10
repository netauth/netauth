package ctl2

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/NetAuth/NetAuth/pkg/client"
)

var (
	entityDestroyCmd = &cobra.Command{
		Use:   "destroy <ID>",
		Short: "Destroy an existing entity",
		Long:  entityDestroyLongDocs,
		Args:  cobra.ExactArgs(1),
		Run:   entityDestroyRun,
	}

	entityDestroyLongDocs = `
Destroy the entity with the specified ID.  The entity is deleted
immediately and without confirmation, please ensure you have typed the
ID correctly.

It is possible to remove the entity running the command, but this is
not recommended and may leave your system without any administrative
users.

The caller must posess the DESTROY_ENTITY capability or be a
GLOBAL_ROOT operator for this command to succeed.`

	entityDestroyExample = `$ netauth entity destroy demo
Entity removed successfully`
)

func init() {
	entityCmd.AddCommand(entityDestroyCmd)
}

func entityDestroyRun(cmd *cobra.Command, args []string) {
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

	result, err := c.RemoveEntity(args[0], t)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(result.GetMsg())
}
