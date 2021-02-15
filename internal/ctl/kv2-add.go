package ctl

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/netauth/netauth/pkg/netauth"
)

var (
	kv2AddCmd = &cobra.Command{
		Use:     "add <entity|group> <target> <key> <value>",
		Short:   "Add a single key and value",
		Long:    kv2AddLongDocs,
		Example: kv2AddExample,
		Args:    kv2AddArgs,
		Run:     kv2AddRun,
	}

	kv2AddLongDocs = `
The Add command allows you to add values to a single key that does not
presently exist on either a group or an entity.  Values will be added
in the order you provide, and ordering will be preserved.
`

	kv2AddExample = `
$ netauth kv add entity example key1 value1
$ netauth kv add entity example cosine:phone "1 (555) 867-5309" "1 (555) 888-8888" "1 (555) 090-0461"

$ netauth kv add group example somenamespace:somekey lots of ordered values
`
)

func init() {
	kv2Cmd.AddCommand(kv2AddCmd)
}

func kv2AddArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 4 {
		return fmt.Errorf("this command requires at least 4 arguments")
	}

	tgt := strings.ToUpper(args[0])
	if tgt != "ENTITY" && tgt != "GROUP" {
		return fmt.Errorf("target must be either an entity or a group")
	}
	return nil
}

func kv2AddRun(cmd *cobra.Command, args []string) {
	var err error

	ctx = netauth.Authorize(ctx, token())

	switch strings.ToLower(args[0]) {
	case "entity":
		err = rpc.EntityKVAdd(ctx, args[1], args[2], args[3:])
	case "group":
		err = nil
	}

	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
