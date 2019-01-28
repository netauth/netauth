package ctl

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/NetAuth/NetAuth/pkg/client"
)

var (
	groupExpansionCmd = &cobra.Command{
		Use:     "expansion <group> <INCLUDE|EXCLUDE|DROP> <target>",
		Short:   "Alter group expansions",
		Long:    groupExpansionLongDocs,
		Example: groupExpansionExample,
		Args:    groupExpansionArgs,
		Run:     groupExpansionRun,
	}

	groupExpansionLongDocs = `
The expansion command manages expansion rules for groups.  Expansions
can be a powerful tool to make your server's memberships easier to
manage, but care should be taken to ensure your expansions remain
maintainable.  The expansions system will ensure that cycles are not
introduced to the membership graph, but no checks are performed for
the sanity of the rules requested or the maintainability of the
resulting graph.  Rules require the membership tree to be parsed for
rules at all levels, and use of expansion rules should be carefully
weighed against the performance requirements of your organization.

There are two types of expansions in NetAuth: INCLUDE and EXCLUDE.
Both of these expansions take a target to act on and are applied to a
single group.  In writing, group expansions should be formatted as
<RULE>:target.  For example INCLUDE:sub-group.

The INCLUDE expansion does exactly what the name implies.  Members of
the target group gain membership in the named group without being
added to it directly.  This expansion is convenient for building up
organizational trees where you might want to translate some easily
statable relation into a group membership.  For example the group
"eng" might include all members of "dev" and "ops".  By adding these
exansions the membership of "eng" is kept up to date without
additional effort.

The EXCLUDE expansion is slightly more complicated.  Members of the
target group are excluded from membership in the source group even if
they are otherwise directly members.  This can be useful if you have a
need to prune out some memberships without removing groups from
individuals.  For example if you have contractors that can't access
production data but otherwise need to be members of groups that grant
such access, you could create a new group "production-data" that gates
this access and has an expansion of EXCLUDE:contractors where
"contractors" contains all contractor owned users (possibly even via
includes).  This would allow you to maintain groups that make sense to
humans while still removing people from groups they shouldn't
logically be in.

Removing an expansion can be done by adding an expansion of the DROP
type.  DROP expansions aren't actually expansions, but they select
existing rules to remove.`

	groupExpansionExample = `$ netauth group expansion example-group include example-group2
Nesting updated successfully
`
)

func init() {
	groupCmd.AddCommand(groupExpansionCmd)
}

func groupExpansionArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 3 {
		return fmt.Errorf("This command takes exactly 3 arguments")
	}

	m := strings.ToUpper(args[1])
	if m != "INCLUDE" && m != "EXCLUDE" && m != "DROP" {
		return fmt.Errorf("mode must be one of INCLUDE, EXCLUDE, or DROP")
	}
	return nil
}

func groupExpansionRun(cmd *cobra.Command, args []string) {
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

	// Apply the expansion
	result, err := c.ModifyGroupExpansions(t, args[0], args[2], args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(result.GetMsg())
}
