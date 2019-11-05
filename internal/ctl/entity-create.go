package ctl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/netauth/netauth/pkg/netauth"
)

var (
	newEntityID string
	newNumber   int
	newSecret   string

	entityCreateCmd = &cobra.Command{
		Use:   "create <ID>",
		Short: "Create a new entity with the specified ID",
		Long:  entityCreateLongDocs,
		Args:  cobra.ExactArgs(1),
		Run:   entityCreateRun,
	}

	entityCreateLongDocs = `
Create an entity with the specified ID.  Though there are no strict
requirements on the ID beyond it being a single word that is globally
unique, it is strongly encouraged to make it exclusively of lower case
letters and numbers.  For the best compatibility, it is recommended to
start with a letter only.

Additional fields can be specified on the command line such as the
number to assign or the initial secret to set.  If left blank the
number will be chosen as the next unassigned number, and the secret
will be prompted for.  To create an entity with an unset secret,
specify the empty string as the initial secret.

The caller must posess the CREATE_ENTITY capability or be a
GLOBAL_ROOT operator for this command to succeed.`

	entityCreateExample = `$ netauth entity create demo
Initial Secret for demo:
New entity created successfully`
)

func init() {
	entityCmd.AddCommand(entityCreateCmd)
	entityCreateCmd.Flags().IntVar(&newNumber, "number", -1, "Number to assign.")
	entityCreateCmd.Flags().StringVar(&newSecret, "initial-secret", "", "Initial secret.")
}

func entityCreateRun(cmd *cobra.Command, args []string) {
	newEntityID = args[0]
	if newSecret == "" {
		newSecret = getSecret(fmt.Sprintf("Initial Secret for %s: ", newEntityID))
	}

	ctx = netauth.Authorize(ctx, token())

	if err := rpc.EntityCreate(ctx, newEntityID, newSecret, newNumber); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Entity Created")
}
