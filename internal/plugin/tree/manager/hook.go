package manager

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"

	"github.com/netauth/netauth/internal/plugin/tree/common"

	pb "github.com/netauth/protocol"
)

// EntityHook satisfies the type for the entity tree system.
type EntityHook struct {
	action common.PluginAction
	mref   *Manager
}

// Name returns the dynamically generated name based on which plugin
// action this hook will invoke.
func (h EntityHook) Name() string {
	return string(fmt.Sprintf("plugin-%s", strings.ToLower(h.action.String())))
}

// Priority looks up the hook priority from the action it is
// performing.
func (h EntityHook) Priority() int {
	return common.AutoHookPriority[h.action]
}

// Run invokes each registered plugin in a non-deterministic order.
// The only order that is guaranteed by this interface is that the
// actions will be called in the same place in the chain each time.
func (h EntityHook) Run(e, de *pb.Entity) error {
	opts := common.PluginOpts{
		Action:     h.action,
		Entity:     e,
		DataEntity: de,
	}

	res, err := h.mref.InvokeEntityProcessing(opts)
	if err != nil {
		return err
	}

	// Danger will robinson!  This Reset makes this merge behave
	// like a first time load event which is what is needed to
	// allow the plugins to write to parts of the entities.
	e.Reset()
	proto.Merge(e, &res.Entity)

	return nil
}

// GroupHook satisfies the type for the group tree system.
type GroupHook struct {
	action common.PluginAction
	mref   *Manager
}

// Name returns the dynamically generated name based on which plugin
// action this hook will invoke.
func (h GroupHook) Name() string {
	return string(fmt.Sprintf("plugin-%s", strings.ToLower(h.action.String())))
}

// Priority looks up the hook priority from the action it is
// performing.
func (h GroupHook) Priority() int {
	return common.AutoHookPriority[h.action]
}

// Run invokes each registered plugin in a non-deterministic order.
// The only order that is guaranteed by this interface is that the
// actions will be called in the same place in the chain each time.
func (h GroupHook) Run(g, dg *pb.Group) error {
	opts := common.PluginOpts{
		Action:    h.action,
		Group:     g,
		DataGroup: dg,
	}

	res, err := h.mref.InvokeGroupProcessing(opts)
	if err != nil {
		return err
	}

	g.Reset()
	proto.Merge(g, &res.Group)

	return nil
}
