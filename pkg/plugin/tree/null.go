package tree

import (
	"context"

	pb "github.com/netauth/protocol"
)

// NullPlugin represents a plugin with no behavior at all.  The intent
// is that it can be embedded into real plugins to stub out the
// functions that the plugin does not wish to implement.
type NullPlugin struct{}

// EntityCreate is called after an entity has been fully created, but
// before the entity has been written to disk.  This hook is an
// excellent time to create entities in other systems.
func (NullPlugin) EntityCreate(_ context.Context, e, de pb.Entity) (pb.Entity, error) {
	return e, nil
}

// EntityUpdate is called after an entity update has occurred, but
// before the updated entity has been written to disk.  This call is
// the appropriate time to update metadata in other systems if it has
// been changed.
func (NullPlugin) EntityUpdate(_ context.Context, e pb.Entity) (pb.Entity, error) {
	return e, nil
}

// EntityLock is called after the lock flag has been set on an entity,
// but before the updated entity has been written to disk.  This call
// is the appropriate point to lock entities in other systems if locks
// are propogated.
func (NullPlugin) EntityLock(_ context.Context, e pb.Entity) (pb.Entity, error) {
	return e, nil
}

// EntityUnlock is called after the lock flag has been cleared for an
// entity.  This call should be used to propagate entity unlocks to
// other systems.
func (NullPlugin) EntityUnlock(_ context.Context, e pb.Entity) (pb.Entity, error) {
	return e, nil
}

// EntityDestroy is called when an entity is being completely removed
// from the system.  At this point in the call chain the entity has
// not been fully removed, but removal will continue if no errors are
// encountered in the processing chain.
func (NullPlugin) EntityDestroy(_ context.Context, e pb.Entity) (pb.Entity, error) {
	return e, nil
}

// GroupCreate is called when a group has been created, but not yet
// written to the storage backend.  This call can be used to create
// matching groups in remote systems.
func (NullPlugin) GroupCreate(_ context.Context, g pb.Group) (pb.Group, error) {
	return g, nil
}

// GroupUpdate is called when the metadata on a group has been
// updated.  If propogating group memberships to another system, a
// periodic approach is recommended, as that will be much cleaner than
// trying to watch for updates and trigger scans dynamically.
func (NullPlugin) GroupUpdate(_ context.Context, g pb.Group) (pb.Group, error) {
	return g, nil
}

// GroupDestroy is called while a group is being fully removed from
// the server.  Groups should never be fully removed, but if they are
// to be destroyed then this function will allow you to propagate this
// destruction to other systems.
func (NullPlugin) GroupDestroy(_ context.Context, g pb.Group) (pb.Group, error) {
	return g, nil
}

// PreSecretChange change is called before an entity has completed a
// secret change, but after the request has been authenticated.
func (NullPlugin) PreSecretChange(_ context.Context, e, de pb.Entity) (pb.Entity, error) {
	return e, nil
}

// PostSecretChange is called after a secret change has occurred, but
// before it is committed to permanent storage.
func (NullPlugin) PostSecretChange(_ context.Context, e, de pb.Entity) (pb.Entity, error) {
	return e, nil
}

// PreAuthCheck is called before an entity has successfully
// authenticated.
func (NullPlugin) PreAuthCheck(_ context.Context, e, de pb.Entity) (pb.Entity, error) {
	return e, nil
}

// PostAuthCheck is called after an entity has successfully
// authenticated.  This method will not be called in the event of an
// authentication failure.
func (NullPlugin) PostAuthCheck(_ context.Context, e, de pb.Entity) (pb.Entity, error) {
	return e, nil
}
