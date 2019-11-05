package ctl

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/netauth/netauth/pkg/netauth"
)

var (
	groupKVCmd = &cobra.Command{
		Use:     "kv <group> <UPSERT|CLEARFUZZY|CLEAREXACT|READ> <key> [value]",
		Short:   "Manage KV storage on an group",
		Long:    groupKVLongDocs,
		Example: groupKVExample,
		Args:    kvArgs,
		Run:     groupKVRun,
	}

	groupKVLongDocs = `
The KV subsystem allows NetAuth to store additional arbitrary
metadata.  Use of this system should be carefully balanced against the
performance impact since this data is stored on groups directly, and
as such can impact access times.

The KV system supports indexed keys, which are of the form key{index}
and are sortable by the client.  For example, if you had multiple
phone numbers that you wanted to keep in order based on the order in
which they are preferred.  The following arrangement would accomplish
this ordering:

	phone{0}: 1 (555) 867-5309
	phone{1}: 1 (555) 888-8888
	phone{2}: 1 (555) 090-0461

If you wanted to change a single key, you could either upsert it which
will insert or update as necessary, or you could remove it.  To remove
the key use either CLEARFUZZY or CLEAREXACT.  The exact variant allows
you to specify the exact key with index to clear, whereas the fuzzy
version doesn't check the index before clearing (useful for bulk
removing a key).
`

	groupKVExample = `$ netauth group kv demo2 upsert phone{0} "1 (555) 867-5309"
$ netauth group kv demo2 upsert phone{1} "1(555) 888-8888"
$ netauth group kv demo2 upsert phone{2} "1(555) 090-0461"

$ netauth group kv demo2 read phone
phone{0}: 1 (555) 867-5309
phone{1}: 1 (555) 888-8888
phone{2}: 1 (555) 090-0461

$ netauth group kv demo2 clearexact phone{1}
$ netauth group kv demo2 read phone
phone{0}: 1 (555) 867-5309
phone{2}: 1 (555) 090-0461

$ netauth group kv demo2 clearfuzzy phone
$ neatuth group kv demo2 read phone
`
)

func init() {
	groupCmd.AddCommand(groupKVCmd)
}

func groupKVRun(cmd *cobra.Command, args []string) {
	// Parse arguments
	action := strings.ToUpper(args[1])
	key := args[2]
	val := ""
	if action == "UPSERT" {
		val = args[3]
	}

	// Get the authorization token if needed
	if action != "READ" {
		ctx = netauth.Authorize(ctx, token())
	}
	// Query the server
	result, err := rpc.GroupUM(ctx, args[0], args[1], key, val)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for k, v := range result {
		fmt.Printf("%s:\n", k)
		for _, s := range v {
			fmt.Printf("  %s\n", s)
		}
	}
}
