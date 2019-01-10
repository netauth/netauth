package ctl2

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/NetAuth/NetAuth/pkg/client"
)

var (
	entityInfoFields string

	entityInfoCmd = &cobra.Command{
		Use:   "info <entity>",
		Short: "Fetch information on an existing entity",
		Long:  entityInfoLongDocs,
		Args:  cobra.ExactArgs(1),
		Run:   entityInfoRun,
	}

	entityInfoLongDocs = `
The info command can return information on any entity known to the
server.  The output may be filtered with the --fields option which
takes a comma seperated list of field names to display.  `

	entityInfoExample = `$ netauth entity info demo2
ID: demo2
Number: 9

$ netauth entity info --fields ID demo2
ID: demo2`
)

func init() {
	entityCmd.AddCommand(entityInfoCmd)
	entityInfoCmd.Flags().StringVar(&entityInfoFields, "fields", "", "Fields to be displayed")
}

func entityInfoRun(cmd *cobra.Command, args []string) {
	// Grab a client
	c, err := client.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Obtain entity info
	entity, err := c.EntityInfo(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Print the fields
	printEntity(entity, entityInfoFields)
}
