package ctl

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/netauth/netauth/pkg/netauth"
)

var (
	kv2ReplaceCmd = &cobra.Command{
		Use:     "replace <entity|group> <target> <key> <value>",
		Short:   "Replace a single key and value",
		Long:    kv2ReplaceLongDocs,
		Example: kv2ReplaceExample,
		Args:    kv2ReplaceArgs,
		Run:     kv2ReplaceRun,
	}

	kv2ReplaceLongDocs = `
The replace command allows you to overwrite the values for a single
key that already exists on an entity or group.  It is identical to add
with the exception that the key must already exist.
`

	kv2ReplaceExample = `
$ netauth kv add entity example key1 value1
$ netauth kv replace entity example key1 value2 value3
`
)

func init() {
	kv2Cmd.AddCommand(kv2ReplaceCmd)
}

func kv2ReplaceArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 4 {
		return fmt.Errorf("this command requires at least 4 arguments")
	}

	tgt := strings.ToUpper(args[0])
	if tgt != "ENTITY" && tgt != "GROUP" {
		return fmt.Errorf("target must be either an entity or a group")
	}
	return nil
}

func kv2ReplaceRun(cmd *cobra.Command, args []string) {
	var err error

	ctx = netauth.Authorize(ctx, token())

	switch strings.ToLower(args[0]) {
	case "entity":
		err = rpc.EntityKVReplace(ctx, args[1], args[2], args[3:])
	case "group":
		err = rpc.GroupKVReplace(ctx, args[1], args[2], args[3:])
	}

	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
