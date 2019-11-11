package ctl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	groupInfoFields string

	groupInfoCmd = &cobra.Command{
		Use:     "info <group>",
		Short:   "Fetch information on an existing group",
		Long:    groupInfoLongDocs,
		Example: groupInfoExample,
		Args:    cobra.ExactArgs(1),
		Run:     groupInfoRun,
	}

	groupInfoLongDocs = `
The info command retursn information on any group known to the server.
The output may be filtered with the --fields option which takes a
comma separated list of field names to display.`

	groupInfoExample = `$ netauth group info example-group
Name: example-group
Display Name:
Number: 10
Expansion: INCLUDE:example-group2`
)

func init() {
	groupCmd.AddCommand(groupInfoCmd)
	groupInfoCmd.Flags().StringVar(&groupInfoFields, "fields", "", "Fields to be displayed")
}

func groupInfoRun(cmd *cobra.Command, args []string) {

	// Obtain group info
	result, sub, err := rpc.GroupInfo(ctx, args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Print the fields
	printGroup(result, groupInfoFields)

	if len(sub) > 0 {
		fmt.Println("The following groups are managed by this group:")
		for _, g := range sub {
			fmt.Printf("  - %s\n", g.GetName())
		}
	}
}
