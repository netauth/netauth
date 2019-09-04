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
func Main() {
	var pluginMap = map[string]plugin.Plugin{
		"treep": &common.GoPluginRPC{},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
