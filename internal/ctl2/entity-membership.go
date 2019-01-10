package ctl2

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/NetAuth/NetAuth/pkg/client"

	pb "github.com/NetAuth/Protocol"
)

var (
	entityMembershipCmd = &cobra.Command{
		Use:     "membership <entity> <ADD|DROP> <group>",
		Short:   "Add or remove direct group memberships",
		Long:    entityMembershipLongDocs,
		Example: entityMembershipExample,
		Args:    entityMembershipArgs,
		Run:     entityMembershipRun,
	}

	entityMembershipLongDocs = `
The membership command adds and removes groups from an entity.  These
groups are direct memberships that are only influenced by EXCLUDE
expansions.

The caller must posses the MODIFY_GROUP_MEMBERS capability or be a
member of the group that is listed to manage the membership of the
target group.`

	entityMembershipExample = `$ netauth entity membership demo2 add demo-group
Membership updated successfully

$ netauth entity membership demo2 drop demo-group
Membership updated successfully`
)

func init() {
	entityCmd.AddCommand(entityMembershipCmd)
}

func entityMembershipArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 3 {
		return fmt.Errorf("This command takes exactly 3 arguments")
	}

	m := strings.ToUpper(args[1])
	if m != "ADD" && m != "DROP" {
		return fmt.Errorf("Mode must be one of ADD or DROP")
	}

	return nil
}

func entityMembershipRun(cmd *cobra.Command, args []string) {
	// Grab a client
	c, err := client.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get the authorization token
	t, err := getToken(c, viper.GetString("entity"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	result := &pb.SimpleResult{}
	switch strings.ToUpper(args[1]) {
	case "ADD":
		result, err = c.AddEntityToGroup(t, args[2], args[0])
	case "DROP":
		result, err = c.AddEntityToGroup(t, args[2], args[0])
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(result.GetMsg())
}
