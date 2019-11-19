package ctl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	entityMembershipsFields string

	entityMembershipsCmd = &cobra.Command{
		Use:     "memberships <entity>",
		Short:   "Memberships held by the specified entity",
		Long:    entityMembershipsLongDocs,
		Example: entityMembershipsExample,
		Args:    cobra.ExactArgs(1),
		Run:     entityMembershipsRun,
	}

	entityMembershipsLongDocs = `
The memberships command returns the memberships held by a particular
entity.  By default the output will include all attributes set on any
returned group.  To filter attributes use the --fields command to
specify a comma separated list of groups that you wish to return.
`

	entityMembershipsExample = `$ netauth entity memberships demo2
Name: demo-group
Display Name: Temporary Demo Group
Number: 9

$ netauth entity memberships demo2 --fields DisplayName
Display Name: Temporary Demo Group
`
)

func init() {
	entityCmd.AddCommand(entityMembershipsCmd)
	entityMembershipsCmd.Flags().StringVar(&entityMembershipsFields, "fields", "", "Fields to be displayed")
}

func entityMembershipsRun(cmd *cobra.Command, args []string) {
	res, err := rpc.EntityGroups(ctx, args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Print the fields
	for i, g := range res {
		printGroup(g, entityMembershipsFields)
		if i < len(res)-1 {
			fmt.Println("---")
		}
	}
}
