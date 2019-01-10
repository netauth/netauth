package ctl2

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/NetAuth/NetAuth/pkg/client"
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
comma seperated list of field names to display.`

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
	// Grab a client
	c, err := client.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Obtain group info
	result, err := c.GroupInfo(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Print the fields
	printGroup(result.GetGroup(), groupInfoFields)
	if len(result.GetManaged()) > 0 {
		fmt.Printf("The following group(s) are managed by %s\n", args[0])
	}
	for _, gn := range result.GetManaged() {
		fmt.Printf("  - %s\n", gn)
	}
}
