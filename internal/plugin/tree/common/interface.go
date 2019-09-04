package common

//go:generate stringer -type=PluginAction

import (
	"net/rpc"

	pb "github.com/NetAuth/Protocol"
)

// PluginAction is used to swich handlers inside a ProcessEntity or
// ProcessGroup handler.
type PluginAction int

// These constants are used to switch actions inside the plugins
// themselves.
const (
	EntityCreate PluginAction = iota
	EntityUpdate
	EntityLock
	EntityUnlock
	EntityDestroy

	GroupCreate
	GroupUpdate
	GroupDestroy

	PreSecretChange
	PostSecretChange
	PreAuthCheck
	PostAuthCheck
)

var (
	// AutoEntityActions is a list of all actions which get
	// automatically generated hooks inserted into the tree
	// processing system.
	AutoEntityActions = [...]PluginAction{
		EntityCreate,
		EntityUpdate,
		EntityLock,
		EntityUnlock,
		EntityDestroy,

		PreSecretChange,
		PostSecretChange,
		PreAuthCheck,
		PostAuthCheck,
	}

	// AutoGroupActions is the same as AutoEntityActions, but has
	// been split out since entities and groups have different
	// signatures for their respective hooks.
	AutoGroupActions = [...]PluginAction{
		GroupCreate,
		GroupUpdate,
		GroupDestroy,
	}

	// AutoHookPriority is used to determine where a hook is to be
	// sequenced in to a chain based on priority.
	AutoHookPriority = map[PluginAction]int{
		EntityCreate:  70,
		EntityUpdate:  70,
		EntityLock:    70,
		EntityUnlock:  70,
		EntityDestroy: 70,

		GroupCreate:  70,
		GroupUpdate:  70,
		GroupDestroy: 70,

		PreSecretChange:  40,
		PostSecretChange: 60,
		PreAuthCheck:     15,
		PostAuthCheck:    60,
	}
)

// Plugin is the type for plugins that extend the functionality of the
// built in tree management system.  The most common type of plugin is
// one that will propogate changes to a system external to NetAuth.
type Plugin interface {
	EntityCreate(pb.Entity) (pb.Entity, error)
	EntityUpdate(pb.Entity) (pb.Entity, error)
	EntityLock(pb.Entity) (pb.Entity, error)
	EntityUnlock(pb.Entity) (pb.Entity, error)
	EntityDestroy(pb.Entity) error

	GroupCreate(pb.Group) (pb.Group, error)
	GroupUpdate(pb.Group) (pb.Group, error)
	GroupDestroy(pb.Group) error

	PreSecretChange(pb.Entity, pb.Entity) (pb.Entity, error)
	PostSecretChange(pb.Entity, pb.Entity) (pb.Entity, error)
	PreAuthCheck(pb.Entity, pb.Entity) (pb.Entity, error)
	PostAuthCheck(pb.Entity, pb.Entity) (pb.Entity, error)
}

// GoPlugin is the actual interface that's exposed across the link.
type GoPlugin interface {
	ProcessEntity(PluginOpts, *PluginResult) error
	ProcessGroup(PluginOpts, *PluginResult) error
}

// GoPluginClient is an RPC Servable type that can be used with Hashicorp's
// go-plugin in order to provide the transport for the actual Plugin
// interface.
type GoPluginClient struct {
	GoPlugin
	client *rpc.Client
}

// GoPluginRPC is a binding only type that's used to provide the
// interface required by go-plugin.
type GoPluginRPC struct{}

// PluginOpts provides a clean transport for data that needs to be fed
// into a plugin.  Note that this is used for both group and entity
// operations, but not all fields are required to be populated.
type PluginOpts struct {
	Action     PluginAction
	Entity     *pb.Entity
	DataEntity *pb.Entity
	Group      *pb.Group
	DataGroup  *pb.Group
}

// PluginResult is returned by group and entity operations in plugins
// and provides a container for data to be passed back along the RPC
// connection.
type PluginResult struct {
	Entity pb.Entity
	Group  pb.Group
}

// GoPluginServer implements the net/rpc server that GoPluginRPC
// talks to.
type GoPluginServer struct{}
