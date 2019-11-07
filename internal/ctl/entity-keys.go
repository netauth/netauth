package ctl

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/netauth/netauth/pkg/netauth"
)

var (
	keyEntity string
	keyMode   string
	keyType   string
	key       string

	entityKeysCmd = &cobra.Command{
		Use:     "key",
		Short:   "Manage keys on an entity",
		Long:    entityKeysLongDocs,
		Example: entityKeysExample,
		Args:    keyArgs,
		Run:     entityKeysRun,
	}

	entityKeysLongDocs = `

The keys command manages the keys that are stored directly on an
entity.  Since the metadata for entities is public it is important to
only ever store *public* keys on the entity.  Most commonly this
feature would be used to store SSH keys that should be trusted across
the network.

The default key type is always SSH, and keys are matched exactly.  It
can be useful to copy and paste a key from the list output to remove
it.`

	entityKeysExample = `$ netauth entity key add SSH "ssh-rsa this-is-too-short-but-whatever root@everywhere"
$ netauth entity key list
Type: SSH; Key: ssh-rsa this-is-too-short-but-whatever root@everywhere
$ netauth entity key drop "ssh-rsa this-is-too-short-but-whatever root@everywhere"
$ netauth entity key list
`
)

func init() {
	entityCmd.AddCommand(entityKeysCmd)
	entityKeysCmd.Flags().StringVar(&keyEntity, "entityID", "", "Entity to change keys for (omit for request entity)")
}

func keyArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("you must specify at minimum a mode")
	} else if len(args) > 3 {
		return fmt.Errorf("too many arguments, is your key quoted?")
	}

	m := strings.ToUpper(args[0])
	if m == "ADD" && len(args) != 3 {
		return fmt.Errorf("ADD requires a keyType and key")
	} else if m == "DROP" && len(args) != 2 {
		return fmt.Errorf("DROP requires at most a key")
	} else if m == "READ" && len(args) > 2 {
		return fmt.Errorf("READ takes at most one argument")
	} else if m != "ADD" && m != "DROP" && m != "READ" {
		return fmt.Errorf("mode must be one of ADD, DEL, or LIST")
	}
	return nil
}

func entityKeysRun(cmd *cobra.Command, args []string) {
	if keyEntity == "" {
		keyEntity = viper.GetString("entity")
	}

	keyMode = strings.ToUpper(args[0])
	switch keyMode {
	case "ADD":
		keyType = args[1]
		key = args[2]
	case "DROP":
		key = args[1]
	case "READ":
		if len(args) == 2 {
			keyType = strings.ToUpper(args[1])
		} else {
			keyType = "SSH"
		}
	}

	if keyMode != "LIST" {
		ctx = netauth.Authorize(ctx, token())
	}

	res, err := rpc.EntityKeys(ctx, keyEntity, keyMode, keyType, key)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for keyType, keys := range res {
		fmt.Printf("%s:\n", keyType)
		for _, key := range keys {
			fmt.Printf("  %s\n", key)
		}
	}
}
