package ctl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	entitySearchFields string

	entitySearchCmd = &cobra.Command{
		Use:     "search <expression>",
		Short:   "Search entities on the server",
		Long:    entitySearchLongDocs,
		Example: entitySearchExample,
		Args:    cobra.ExactArgs(1),
		Run:     entitySearchRun,
	}

	entitySearchLongDocs = `
The search command allows complex searching within entities.  This
command takes a single argument which is the search expression, be
sure to quote the expression if making a complex query.

All set fields on returned entities will be displayed.  To display
only certain fields pass a comma seperated list to the --fields
argument of the field names you wish to display.

Some fields on entities are part of the metadata, to address these
fields in a search prefix them with 'meta.' as in 'meta.DisplayName'.`

	entitySearchExample = `$ netauth entity search 'ID:demo*'
ID: demo2
Number: 9
ID: demo3
Number: 10
ID: demo4
Number: 11

$ netauth entity search 'meta.Shell: /bin/bash'
ID: demo3
Number: 10
shell: /bin/bash
`
)

func init() {
	entityCmd.AddCommand(entitySearchCmd)
	entitySearchCmd.Flags().StringVar(&entitySearchFields, "fields", "", "Fields to be displayed")
}

func entitySearchRun(cmd *cobra.Command, args []string) {
	// Obtain entity info
	res, err := rpc.EntitySearch(ctx, args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Print the fields
	for _, e := range res {
		printEntity(e, entitySearchFields)
	}
}
