package manager

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/internal/plugin/tree/common"
	"github.com/NetAuth/NetAuth/internal/tree/util"

	pb "github.com/NetAuth/Protocol"
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
// The only order that is gauranteed by this interface is that the
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

	proto.Merge(e, &res.Entity)

	// This step is silly, but is needed to ensure that
	// proto.Merge hasn't done something silly above.
	if e.Meta != nil {
		e.Meta.Groups = util.DedupStringSlice(e.Meta.Groups)
		e.Meta.Keys = util.DedupStringSlice(e.Meta.Keys)
		e.Meta.UntypedMeta = util.DedupStringSlice(e.Meta.UntypedMeta)
		e.Meta.Capabilities = util.DedupCapabilitySlice(e.Meta.Capabilities)
	}

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
// The only order that is gauranteed by this interface is that the
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

	proto.Merge(g, &res.Group)

	// This step is silly, but is needed to ensure that
	// proto.Merge hasn't done something silly above.
	g.Expansions = util.DedupStringSlice(g.Expansions)
	g.UntypedMeta = util.DedupStringSlice(g.UntypedMeta)
	g.Capabilities = util.DedupCapabilitySlice(g.Capabilities)

	return nil
}
