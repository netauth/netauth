package tree

import (
	"github.com/NetAuth/NetAuth/internal/plugin/tree/provider"
	"github.com/NetAuth/NetAuth/internal/plugin/tree/common"
)

// Plugin is an alias to an internal type that all tree modifying
// plugins must satisfy.
type Plugin = common.Plugin

// PluginMain is called with an interface to serve as the plugin.
// This function never returns.
func PluginMain(i Plugin) {
	provider.Main(i)
}
