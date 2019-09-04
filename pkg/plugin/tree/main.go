package tree

import (
	"github.com/NetAuth/NetAuth/internal/plugin/tree/provider"
	"github.com/NetAuth/NetAuth/internal/plugin/tree/common"
)

// PluginMain is called with an interface to serve as the plugin.
// This function never returns.
func PluginMain(i common.Plugin) {
	provider.Main(i)
}
