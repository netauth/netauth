package common

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

func (p *GoPluginServer) Server(*plugin.MuxBroker) (interface{}, error) {
	return &GoPluginServer{}, nil
}

func (GoPluginServer) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &GoPluginClient{client: c}, nil
}

func (p *GoPluginServer) ProcessEntity(opts PluginOpts, res *PluginResult) error {
	res.Entity = opts.Entity
	return nil
}

func (p *GoPluginServer) ProcessGroup(opts PluginOpts, res *PluginResult) error {
	res.Group = opts.Group
	return nil
}

func (c *GoPluginClient) ProcessEntity(opts PluginOpts) (PluginResult, error) {
	var res PluginResult
	err := c.client.Call("Plugin.ProcessEntity", opts, &res)
	return res, err
}

func (c *GoPluginClient) ProcessGroup(opts PluginOpts) (PluginResult, error) {
	var res PluginResult
	err := c.client.Call("Plugin.ProcessGroup", opts, &res)
	return res, err
}
