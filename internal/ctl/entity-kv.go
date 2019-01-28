package ctl

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/NetAuth/NetAuth/pkg/client"
)

var (
	entityKVCmd = &cobra.Command{
		Use:     "kv <entity> <UPSERT|CLEARFUZZY|CLEAREXACT|READ> <key> [value]",
		Short:   "Manage KV storage on an entity",
		Long:    entityKVLongDocs,
		Example: entityKVExample,
		Args:    entityKVArgs,
		Run:     entityKVRun,
	}

	entityKVLongDocs = `
The KV subsystem allows NetAuth to store additional arbitrary
metadata.  Use of this system should be carefully balanced against the
performance impact since this data is stored on entities directly, and
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

	entityKVExample = `$ netauth entity kv demo2 upsert phone{0} "1 (555) 867-5309"
$ netauth entity kv demo2 upsert phone{1} "1(555) 888-8888"
$ netauth entity kv demo2 upsert phone{2} "1(555) 090-0461"

$ netauth entity kv demo2 read phone
phone{0}: 1 (555) 867-5309
phone{1}: 1 (555) 888-8888
phone{2}: 1 (555) 090-0461

$ netauth entity kv demo2 clearexact phone{1}
$ netauth entity kv demo2 read phone
phone{0}: 1 (555) 867-5309
phone{2}: 1 (555) 090-0461

$ netauth entity kv demo2 clearfuzzy phone
$ neatuth entity kv demo2 read phone
`
)

func init() {
	entityCmd.AddCommand(entityKVCmd)
}

func entityKVArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("This command takes at least 3 arguments")
	}
	action := strings.ToUpper(args[1])
	if action == "UPSERT" && len(args) != 4 {
		return fmt.Errorf("Upsert requires a key and a value")
	}

	switch action {
	case "UPSERT":
		return nil
	case "CLEARFUZZY":
		return nil
	case "CLEAREXACT":
		return nil
	case "READ":
		return nil
	default:
		return fmt.Errorf("Action must be one of UPSERT, CLEARFUZZY, CLEAREXACT, or READ")
	}
}

func entityKVRun(cmd *cobra.Command, args []string) {
	// Parse arguments
	action := strings.ToUpper(args[1])
	key := args[2]
	val := ""
	if action == "UPSERT" {
		val = args[3]
	}

	// Grab a client
	c, err := client.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get the authorization token if needed
	t := ""
	if action != "READ" {
		t, err = getToken(c, viper.GetString("entity"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	// Query the server
	result, err := c.ModifyUntypedEntityMeta(t, args[0], args[1], key, val)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	out := []string{}
	for k, v := range result {
		out = append(out, fmt.Sprintf("%s: %s", k, v))
	}
	sort.Strings(out)

	for _, l := range out {
		fmt.Println(l)
	}
}
