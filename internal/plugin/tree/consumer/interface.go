package consumer

import (
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/netauth/netauth/internal/plugin/tree/common"
)

// Ref is a reference to a specific plugin.
type Ref struct {
	*common.GoPluginClient

	path   string
	client plugin.ClientProtocol
	cfg    *plugin.Client
	log    hclog.Logger
}
