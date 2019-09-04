package provider

import (
	"github.com/hashicorp/go-plugin"

	"github.com/NetAuth/NetAuth/internal/plugin/tree/common"
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
func (m mux) HandleEntity(o common.PluginOpts) (common.PluginResult, error) {
	var res common.PluginResult
	var err error

	e := *o.Entity
	de := *o.DataEntity

	switch o.Action {
	case common.EntityCreate:
		res.Entity, err = m.impl.EntityCreate(e)
	case common.EntityUpdate:
		res.Entity, err = m.impl.EntityUpdate(e)
	case common.EntityLock:
		res.Entity, err = m.impl.EntityLock(e)
	case common.EntityUnlock:
		res.Entity, err = m.impl.EntityUnlock(e)
	case common.EntityDestroy:
		err = m.impl.EntityDestroy(e)
	case common.PreSecretChange:
		res.Entity, err = m.impl.PreSecretChange(e, de)
	case common.PostSecretChange:
		res.Entity, err = m.impl.PostSecretChange(e, de)
	case common.PreAuthCheck:
		res.Entity, err = m.impl.PreAuthCheck(e, de)
	case common.PostAuthCheck:
		res.Entity, err = m.impl.PostAuthCheck(e, de)
	default:
		res.Entity = e
		err = nil
	}
	return res, err
}

// HandleGroup is the same as HandleEntity, but acts on group
// functions.  These methods are split in order to keep the context of
// the switch statements to a reasonable size.
func (m mux) HandleGroup(o common.PluginOpts) (common.PluginResult, error) {
	var res common.PluginResult
	var err error

	g := *o.Group

	switch o.Action {
	case common.GroupCreate:
		res.Group, err = m.impl.GroupCreate(g)
	case common.GroupUpdate:
		res.Group, err = m.impl.GroupUpdate(g)
	case common.GroupDestroy:
		err = m.impl.GroupDestroy(g)
	default:
		res.Group = g
		err = nil
	}
	return res, err
}
