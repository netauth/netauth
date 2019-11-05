package ctl

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/netauth/netauth/pkg/netauth"

	pb "github.com/NetAuth/Protocol"
)

var (
	direct bool

	systemCapabilitiesCmd = &cobra.Command{
		Use:     "capability <identifier> <ADD|DEL> <capability>",
		Short:   "Manage internal system capabilities",
		Long:    systemCapabilityLongDocs,
		Example: systemCapabilityExample,
		Args:    systemCapabilityArgs,
		Run:     systemCapabilitiesRun,
	}

	systemCapabilityLongDocs = `
NetAuth makes use of a capabilities based system for internal access
control.  The capabilities command can add and remove capabilities
from entities and groups.  The preferred mechanism for access control
should always be to gain capabilities by being in a group that has
them, rather than having access applied to entities directly.  A
description of each capability follows:

  GLOBAL_ROOT - Confers all other capabilities implicitly.  This power
    is used to bootstrap the server and should be reserved to super
    administrators that would otherwise be able to obtain this power.

  CREATE_ENTITY - Allow the creation of entities.

  DESTROY_ENTITY - Allows the destruction of entities.

  MODIFY_ENTITY_META - Allows modification of entity metadata.

  MODIFY_ENTITY_KEYS - Allows modification of entity public keys.
    Entities are able to change their own keys without this capability.

  CHANGE_ENTITY_SECRET - Allows modification of entity secrest.
    Entities are able to change their own secrets without this
    capability.

  LOCK_ENTITY - Allows setting an entity lock.  Locked entities cannot
    successfully authenticate, even with a correct secret.

  UNLOCK_ENTITY - Allows unlocking an entity.

  CREATE_GROUP - Allows creation of groups.

  DESTROY_GROUP - Allows destruction of groups.

  MODIFY_GROUP_META - Allows the modification of group level metadata.
    This should generally be assigned in conjunction with.

  MODIFY_GROUP_MEMBERS - Allows the modification of group memberships.
    This capability is not needed if the requesting entity is a member
    of a groups designated management group.
`

	systemCapabilityExample = `$ netauth system capability example-group add MODIFY_GROUP_META
Capability Modified

$ netauth system capability --direct demo2 add MODIFY_GROUP_META
You are attempting to add a capability directly to an entity.  This is discouraged!
Capability Modified`
)

func init() {
	systemCmd.AddCommand(systemCapabilitiesCmd)
	systemCapabilitiesCmd.Flags().BoolVar(&direct, "direct", false, "Provided identifier is an entity (discouraged)")
}

func systemCapabilityArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 3 {
		return fmt.Errorf("this command takes 3 arguments")
	}

	if strings.ToUpper(args[1]) != "ADD" && strings.ToUpper(args[1]) != "DROP" {
		return fmt.Errorf("mode must be either ADD or DROP")
	}

	c := strings.ToUpper(args[2])
	if _, ok := pb.Capability_value[c]; !ok {
		return fmt.Errorf("%s is not a known capability", args[2])
	}

	return nil
}

func systemCapabilitiesRun(cmd *cobra.Command, args []string) {
	ctx = netauth.Authorize(ctx, token())

	if err := rpc.SystemCapabilities(ctx, args[0], args[1], args[2], direct); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Capabilities Updated")
}
