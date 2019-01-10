package ctl2

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/NetAuth/NetAuth/pkg/client"
)

var (
	entityMembershipsFields    string
	entityMembershipsIndirects bool

	entityMembershipsCmd = &cobra.Command{
		Use:   "memberships <entity>",
		Short: "Memberships held by the specified entity",
		Long:  entityMembershipLongDocs,
		Example: entityMembershipsExample,
		Args:  cobra.ExactArgs(1),
		Run:   entityMembershipsRun,
	}

	entityMembershipsLongDocs = `
The memberships command returns the memberships held by a particular
entity.  By default the output will include all attributes set on any
returned group.  To filter attributes use the --fields command to
specify a comma seperated list of groups that you wish to return.

The listing will include by default all memberships, including those
gained by group expansions to other groups.  To suppress group
indirects use the option --indirect=false.  Exclude expansions are
processed unconditionally.`

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
	entityMembershipsCmd.Flags().BoolVar(&entityMembershipsIndirects, "indirect", true, "Include indirect memberships")
}

func entityMembershipsRun(cmd *cobra.Command, args []string) {
	// Grab a client
	c, err := client.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	res, err := c.ListGroups(args[0], entityMembershipsIndirects)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Print the fields
	for _, g := range res {
		printGroup(g, entityMembershipsFields)
	}
}
