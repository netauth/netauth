package ctl

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/netauth/netauth/pkg/netauth"
)

var (
	kv2DelCmd = &cobra.Command{
		Use:     "del <entity|group> <target> <key>",
		Short:   "Delete a single key",
		Long:    kv2DelLongDocs,
		Example: kv2DelExample,
		Args:    kv2DelArgs,
		Run:     kv2DelRun,
	}

	kv2DelLongDocs = `
The del command allows you to delete values to a single key that
presently exists on either a group or an entity.
`

	kv2DelExample = `
$ netauth kv del entity example key1
`
)

func init() {
	kv2Cmd.AddCommand(kv2DelCmd)
}

func kv2DelArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 3 {
		return fmt.Errorf("this command requires exactly 3 arguments")
	}

	tgt := strings.ToUpper(args[0])
	if tgt != "ENTITY" && tgt != "GROUP" {
		return fmt.Errorf("target must be either an entity or a group")
	}
	return nil
}

func kv2DelRun(cmd *cobra.Command, args []string) {
	var err error

	ctx = netauth.Authorize(ctx, token())

	switch strings.ToLower(args[0]) {
	case "entity":
		err = rpc.EntityKVDel(ctx, args[1], args[2])
	case "group":
		err = nil
	}

	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
