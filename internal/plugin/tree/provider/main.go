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

func Main() {
	var pluginMap = map[string]plugin.Plugin{
		"treep": &common.GoPluginServer{},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
