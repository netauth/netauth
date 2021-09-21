package provider

import (
	"context"

	"github.com/hashicorp/go-plugin"

	"github.com/netauth/netauth/internal/plugin/tree/common"
)

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "TREE_PLUGIN",
	MagicCookieValue: "treehello",
}

// Main is used to drop into the plugin serving loop that will expose
// this plugin's functions via RPC.
func Main(i common.Plugin) {
	m := mux{impl: i}

	var pluginMap = map[string]plugin.Plugin{
		"treep": &common.GoPluginRPC{Mux: m},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}

type mux struct {
	impl common.Plugin
}

// HandleEntity is a function that provides the selection behavior to
// call the right interface method from the plugin implementation
// based on the action.
func (m mux) HandleEntity(ctx context.Context, o common.PluginOpts) (common.PluginResult, error) {
	var res common.PluginResult
	var err error

	e := *o.Entity
	de := *o.DataEntity

	switch o.Action {
	case common.EntityCreate:
		res.Entity, err = m.impl.EntityCreate(ctx, e, de)
	case common.EntityUpdate:
		res.Entity, err = m.impl.EntityUpdate(ctx, e)
	case common.EntityLock:
		res.Entity, err = m.impl.EntityLock(ctx, e)
	case common.EntityUnlock:
		res.Entity, err = m.impl.EntityUnlock(ctx, e)
	case common.EntityDestroy:
		res.Entity, err = m.impl.EntityDestroy(ctx, e)
	case common.PreSecretChange:
		res.Entity, err = m.impl.PreSecretChange(ctx, e, de)
	case common.PostSecretChange:
		res.Entity, err = m.impl.PostSecretChange(ctx, e, de)
	case common.PreAuthCheck:
		res.Entity, err = m.impl.PreAuthCheck(ctx, e, de)
	case common.PostAuthCheck:
		res.Entity, err = m.impl.PostAuthCheck(ctx, e, de)
	default:
		res.Entity = e
		err = nil
	}
	return res, err
}

// HandleGroup is the same as HandleEntity, but acts on group
// functions.  These methods are split in order to keep the context of
// the switch statements to a reasonable size.
func (m mux) HandleGroup(ctx context.Context, o common.PluginOpts) (common.PluginResult, error) {
	var res common.PluginResult
	var err error

	g := *o.Group

	switch o.Action {
	case common.GroupCreate:
		res.Group, err = m.impl.GroupCreate(ctx, g)
	case common.GroupUpdate:
		res.Group, err = m.impl.GroupUpdate(ctx, g)
	case common.GroupDestroy:
		res.Group, err = m.impl.GroupDestroy(ctx, g)
	default:
		res.Group = g
		err = nil
	}
	return res, err
}
