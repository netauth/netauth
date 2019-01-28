package ctl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/NetAuth/NetAuth/pkg/client"
)

var (
	groupMembersFields string

	groupMembersCmd = &cobra.Command{
		Use:     "members <group>",
		Short:   "Print the members of the specified group",
		Long:    groupMembersLongDocs,
		Example: groupMembersExample,
		Args:    cobra.ExactArgs(1),
		Run:     groupMembersRun,
	}

	groupMembersLongDocs = `
The members command can summon the membership of a particular group.
The output may be filtered with the --fields option which takes a
comma seperated list of fields to be displayed.`

	groupMembersExample = `$ netauth group members example-group
ID: demo2
Number: 9
displayname: Demonstration Entity
ID: demo3
Number: 10
shell: /bin/bash
ID: demo4
Number: 11`
)

func init() {
	groupCmd.AddCommand(groupMembersCmd)
	groupMembersCmd.Flags().StringVar(&groupMembersFields, "fields", "", "Fields to be displayed")
}

func groupMembersRun(cmd *cobra.Command, args []string) {
	// Grab a client
	c, err := client.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	res, err := c.ListGroupMembers(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Print the fields
	for _, e := range res {
		printEntity(e, groupMembersFields)
	}
}
