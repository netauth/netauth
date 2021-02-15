package ctl

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	kv2GetCmd = &cobra.Command{
		Use:     "get <entity|group> <target> <key>",
		Short:   "Retrieve the value of a key",
		Long:    kv2GetLongDocs,
		Example: kv2GetExample,
		Args:    kv2GetArgs,
		Run:     kv2GetRun,
	}

	kv2GetLongDocs = `
The Get command allows you to retrieve the values for a single key
from either an entity or a group.  If an order was provided when the
values were provided to NetAuth, the returned values will be in this
order.
`
	kv2GetExample = `
$ netauth kv get entity example key1
value1
value2
value3
`
)

func init() {
	kv2Cmd.AddCommand(kv2GetCmd)
}

func kv2GetArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 3 {
		return fmt.Errorf("this command requires exactly 3 arguments")
	}

	tgt := strings.ToUpper(args[0])
	if tgt != "ENTITY" && tgt != "GROUP" {
		return fmt.Errorf("target must be either an entity or a group")
	}
	return nil
}

func kv2GetRun(cmd *cobra.Command, args []string) {
	var res []string
	var err error

	switch strings.ToLower(args[0]) {
	case "entity":
		res, err = rpc.EntityKVGet(ctx, args[1], args[2])
	case "group":
		res, err = rpc.GroupKVGet(ctx, args[1], args[2])
	}
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	if len(res) > 0 {
		for i := range res {
			fmt.Println(res[i])
		}
	}
}
