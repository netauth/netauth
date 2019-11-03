package netauth

import (
	"context"
	"errors"
	"fmt"

	pb "github.com/NetAuth/Protocol"
	rpc "github.com/NetAuth/Protocol/v2"
)

// SystemCapabilities handles the modification of capabilities within
// the server.  Capabilities are the core of NetAuth's internal
// permissions system, and allow the holder to perform special actions
// within the server itself.  Capabilities should generally be
// assigned to groups rather than directly to entities, but there are
// valid cases to assign to an entity directly.
func (c *Client) SystemCapabilities(ctx context.Context, target, action, capability string, direct bool) error {
	if err := c.makeWritable(); err != nil {
		return err
	}

	ctx = c.appendMetadata(ctx)
	cap, ok := pb.Capability_value[capability]
	if !ok {
		return fmt.Errorf("%s is not a recognized capability", capability)
	}

	a, ok := rpc.Action_value[action]
	if !ok {
		return errors.New("Action must be one of ADD, DROP")
	}

	r := rpc.CapabilityRequest{
		Direct:     &direct,
		Target:     &target,
		Action:     rpc.Action(a).Enum(),
		Capability: pb.Capability(cap).Enum(),
	}
	_, err := c.rpc.SystemCapabilities(ctx, &r)
	return err
}

// SystemPing pings the server and obtains back a pong if the server
// is healthy.  If the server is not healthy error will be not nil.
// Use this function to gate healthy servers with a load balancer.
func (c *Client) SystemPing(ctx context.Context) error {
	ctx = c.appendMetadata(ctx)
	_, err := c.rpc.SystemPing(ctx, &rpc.Empty{})
	return err
}

// SystemStatus returns detailed status information about the server.
// This information includes a subsystem report and the first failure
// detected during a health check should a failure be detected.
func (c *Client) SystemStatus(ctx context.Context) (*rpc.ServerStatus, error) {
	ctx = c.appendMetadata(ctx)
	return c.rpc.SystemStatus(ctx, &rpc.Empty{})
}
