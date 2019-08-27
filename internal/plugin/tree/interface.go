package tree

import (
	pb "github.com/NetAuth/Protocol"
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
