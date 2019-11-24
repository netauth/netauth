package netauth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"

	pb "github.com/netauth/protocol"
	rpc "github.com/netauth/protocol/v2"
)

// GroupCreate creates a new group with the specified parameters.  If
// you do not require a specific group number you may pass -1 to
// select the next available number.  To make a group managed by
// another group from the start, pass the name of another group here
// as the managed-by value in order to enable delegated management.
func (c *Client) GroupCreate(ctx context.Context, name, displayName, managedBy string, number int) error {
	if err := c.makeWritable(); err != nil {
		return err
	}

	ctx = c.appendMetadata(ctx)
	r := rpc.GroupRequest{
		Group: &pb.Group{
			Name:        &name,
			DisplayName: &displayName,
			ManagedBy:   &managedBy,
			Number:      proto.Int32(int32(number)),
		},
	}
	_, err := c.rpc.GroupCreate(ctx, &r)
	return err
}

// GroupUpdate allows an existing group to be updated.  Only some
// fields on each group can be updated though, so this function will
// silently unset fields that are not permissible to edit.
func (c *Client) GroupUpdate(ctx context.Context, update *pb.Group) error {
	if err := c.makeWritable(); err != nil {
		return err
	}

	ctx = c.appendMetadata(ctx)
	r := rpc.GroupRequest{
		Group: update,
	}

	_, err := c.rpc.GroupUpdate(ctx, &r)
	return err
}

// GroupInfo returns a single group to the caller.  This function does
// not require an authorized context.
func (c *Client) GroupInfo(ctx context.Context, name string) (*pb.Group, []*pb.Group, error) {
	ctx = c.appendMetadata(ctx)
	r := rpc.GroupRequest{
		Group: &pb.Group{
			Name: &name,
		},
	}
	res, err := c.rpc.GroupInfo(ctx, &r)
	if err != nil {
		return nil, nil, err
	}

	res2, err := c.GroupSearch(ctx, fmt.Sprintf("ManagedBy:%s", name))
	if err != nil {
		c.log.Warn("Error searching", "error", err)
		return res.GetGroups()[0], nil, nil
	}

	return res.GetGroups()[0], res2, nil
}

// GroupUM handles operations concerning the untyped key-value store
// on each group.  This data is not directly processed by NetAuth or
// visible in search indexes, but is useful for integrating with 3rd
// party systems as it provides an ideal place to store alternate keys
// or IDs.  Reads may be performed without authentication, writes must
// be authenticated.
func (c *Client) GroupUM(ctx context.Context, target, action, key, value string) (map[string][]string, error) {
	if strings.ToUpper(action) != "READ" {
		if err := c.makeWritable(); err != nil {
			return nil, err
		}
	}

	ctx = c.appendMetadata(ctx)
	action = strings.ToUpper(action)
	a, ok := rpc.Action_value[action]
	if !ok {
		return nil, errors.New("action must be one of UPSERT, CLEARFUZZY, CLEAREXACT, READ")
	}

	r := rpc.KVRequest{
		Target: &target,
		Key:    &key,
		Value:  &value,
		Action: rpc.Action(a).Enum(),
	}

	res, err := c.rpc.GroupUM(ctx, &r)
	if err != nil {
		return nil, err
	}
	if action == "READ" {
		kv := parseKV(res.GetStrings())
		if key != "*" {
			// Asked for a specific key, fish it out and
			// return a much sparser map.
			return map[string][]string{key: kv[key]}, nil
		}
		// Key was *, return everything
		return kv, nil
	}

	// Not in read mode, return nil for both.
	return nil, nil
}

// GroupUpdateRules manages the rules on groups.  These rules can
// transparently include other groups, recursively remove members, or
// reset the behavior of a group to the default.
func (c *Client) GroupUpdateRules(ctx context.Context, group, action, target string) error {
	if err := c.makeWritable(); err != nil {
		return err
	}

	ctx = c.appendMetadata(ctx)
	action = strings.ToUpper(action)

	// Legacy compat shim to allow for the shorter "DROP" to also
	// remove rules.
	if action == "DROP" {
		action = "REMOVE_RULE"
	}

	a, ok := rpc.RuleAction_value[action]
	if !ok {
		return errors.New("action must be one of INCLUDE, EXCLUDE, OR REMOVE_RULE")
	}

	r := rpc.GroupRulesRequest{
		Group: &pb.Group{
			Name: &group,
		},
		Target: &pb.Group{
			Name: &target,
		},
		RuleAction: rpc.RuleAction(a).Enum(),
	}
	_, err := c.rpc.GroupUpdateRules(ctx, &r)
	return err
}

// GroupAddMember adds a member to a group.  Keep in mind that not all
// systems hooking into NetAuth perform synchronous lookups, so
// membership changes may take some time to propagate.
func (c *Client) GroupAddMember(ctx context.Context, group, entity string) error {
	if err := c.makeWritable(); err != nil {
		return err
	}

	ctx = c.appendMetadata(ctx)
	r := rpc.EntityRequest{
		Entity: &pb.Entity{
			ID: &entity,
			Meta: &pb.EntityMeta{
				Groups: []string{group},
			},
		},
	}
	_, err := c.rpc.GroupAddMember(ctx, &r)
	return err
}

// GroupDelMember removes a member from a group.  Keep in mind that
// not all systems hooking into NetAuth perform synchronous lookups,
// so membership changes may take some time to propagate.
func (c *Client) GroupDelMember(ctx context.Context, group, entity string) error {
	if err := c.makeWritable(); err != nil {
		return err
	}

	ctx = c.appendMetadata(ctx)
	r := rpc.EntityRequest{
		Entity: &pb.Entity{
			ID: &entity,
			Meta: &pb.EntityMeta{
				Groups: []string{group},
			},
		},
	}
	_, err := c.rpc.GroupDelMember(ctx, &r)
	return err
}

// GroupDestroy permanently removes a group from the server.  This is
// not recommended as NetAuth does not perform internal referential
// integrity checks, so it is possible to remove a group that has
// rules pointing at it or otherwise create cycles in the graph.  The
// best practices are to keep groups forever.  They're cheap and as
// long as they're not queried they don't represent additional load.
func (c *Client) GroupDestroy(ctx context.Context, name string) error {
	if err := c.makeWritable(); err != nil {
		return err
	}

	ctx = c.appendMetadata(ctx)
	r := rpc.GroupRequest{
		Group: &pb.Group{
			Name: &name,
		},
	}
	_, err := c.rpc.GroupDestroy(ctx, &r)
	return err
}

// GroupMembers returns the membership of a group including any member
// alterations as a result of rules on the group.
func (c *Client) GroupMembers(ctx context.Context, name string) ([]*pb.Entity, error) {
	ctx = c.appendMetadata(ctx)
	r := rpc.GroupRequest{
		Group: &pb.Group{
			Name: &name,
		},
	}
	res, err := c.rpc.GroupMembers(ctx, &r)
	return res.GetEntities(), err
}

// GroupSearch returns a list of groups that satisfy the given search
// expression.  This function requires no authorization.
func (c *Client) GroupSearch(ctx context.Context, expression string) ([]*pb.Group, error) {
	ctx = c.appendMetadata(ctx)
	r := rpc.SearchRequest{
		Expression: &expression,
	}
	res, err := c.rpc.GroupSearch(ctx, &r)
	if err != nil {
		return nil, err

	}
	return res.GetGroups(), err
}
