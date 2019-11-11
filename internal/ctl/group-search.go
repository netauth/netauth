package ctl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	groupSearchFields string

	groupSearchCmd = &cobra.Command{
		Use:     "search <expression>",
		Short:   "Search entities on the server",
		Long:    groupSearchLongDocs,
		Example: groupSearchExample,
		Args:    cobra.ExactArgs(1),
		Run:     groupSearchRun,
	}

	groupSearchLongDocs = `
The search command allows complex searching within groups.  This
command takes a single argument which is the search expression, be
sure to quote the expression if making a complex query.

All set fields on returned groups will be displayed.  To display
only certain fields pass a comma separated list to the --fields
argument of the field names you wish to display.`

	groupSearchExample = `$ netauth group search 'Name:example*'
Name: example-group
Display Name:
Number: 10
Expansion: INCLUDE:example-group2
Name: example-group2
Display Name:
Number: 11
`
)

func init() {
	groupCmd.AddCommand(groupSearchCmd)
	groupSearchCmd.Flags().StringVar(&groupSearchFields, "fields", "", "Fields to be displayed")
}

func groupSearchRun(cmd *cobra.Command, args []string) {
	res, err := rpc.GroupSearch(ctx, args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Print the fields
	for _, g := range res {
		printGroup(g, groupSearchFields)
	}
}
