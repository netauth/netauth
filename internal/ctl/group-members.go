package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/NetAuth/NetAuth/pkg/client"

	"github.com/google/subcommands"
)

type GroupMembersCmd struct {
	ID string
}

func (*GroupMembersCmd) Name() string     { return "group-members" }
func (*GroupMembersCmd) Synopsis() string { return "List members in a named group" }
func (*GroupMembersCmd) Usage() string {
	return `group-members --ID <ID>
List the members of the group identified by <ID>.`
}

func (p *GroupMembersCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.ID, "ID", "", "ID of the group to list")
}

func (p *GroupMembersCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	membersList, err := client.GroupMembers(serverAddr, serverPort, clientID, serviceID, p.ID)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	for _, m := range membersList.Members {
		fmt.Println(m)
	}
	
	return subcommands.ExitSuccess
}
