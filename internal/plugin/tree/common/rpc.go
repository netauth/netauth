package common

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// Server returns a go-plugin compliant interface that handles the
// provider side of the interface.
func (p *GoPluginRPC) Server(*plugin.MuxBroker) (interface{}, error) {
	return &GoPluginServer{Mux: p.Mux}, nil
}

// Client returns a go-plugin compliant interface that handles the
// consumer side of the interface.
func (GoPluginRPC) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &GoPluginClient{client: c}, nil
}

// ProcessEntity on the GoPluginServer type implements a net/rpc
// server method that handles entities.
func (p *GoPluginServer) ProcessEntity(opts PluginOpts, res *PluginResult) error {
	tmp, err := p.Mux.HandleEntity(opts)
	res.Entity = tmp.Entity
	return err
}

// ProcessGroup on the GoPluginServer type implements a net/rpc server
// that handles groups.
func (p *GoPluginServer) ProcessGroup(opts PluginOpts, res *PluginResult) error {
	tmp, err := p.Mux.HandleGroup(opts)
	res.Group = tmp.Group
	return err
}

// ProcessEntity on the GoPluginClient provides a much cleaner
// interface than a raw net/rpc connection.  ProcessEntity handles
// modifications that handle entities only.
func (c *GoPluginClient) ProcessEntity(opts PluginOpts) (PluginResult, error) {
	var res PluginResult
	err := c.client.Call("Plugin.ProcessEntity", opts, &res)
	return res, err
}

// ProcessGroup on the GoPluginClient provides a much cleaner
// interface than a raw net/rpc connection.  ProcessGroup handles
// modifications that handle entities only.
func (c *GoPluginClient) ProcessGroup(opts PluginOpts) (PluginResult, error) {
	var res PluginResult
	err := c.client.Call("Plugin.ProcessGroup", opts, &res)
	return res, err
}
