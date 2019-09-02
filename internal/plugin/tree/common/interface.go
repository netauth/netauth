package common

import (
	"net/rpc"

	pb "github.com/NetAuth/Protocol"
)

type pluginAction int

// These constants are used to switch actions inside the plugins
// themselves.
const (
	EntityCreate pluginAction = iota
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
	PreTokenAuth
	PostTokenAuth
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

	PreSecretChange(pb.Entity) (pb.Entity, error)
	PostSecretChange(pb.Entity) (pb.Entity, error)
	PreAuthCheck(pb.Entity) (pb.Entity, error)
	PostAuthCheck(pb.Entity) (pb.Entity, error)
	PreTokenAuth(pb.Entity) (pb.Entity, error)
	PostTokenAuth(pb.Entity) (pb.Entity, error)
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

type PluginOpts struct {
	Action pluginAction
	Entity pb.Entity
	Group  pb.Group
}

type PluginResult struct {
	Entity pb.Entity
	Group  pb.Group
}

// GoPluginServer implements the net/rpc server that GoPluginRPC
// talks to.
type GoPluginServer struct{}
