package netauth

import (
	"context"
	"errors"
	"sort"
	"strings"

	"google.golang.org/protobuf/proto"

	pb "github.com/netauth/protocol"
	rpc "github.com/netauth/protocol/v2"
)

// EntityCreate creates an entity.  The entity ID must be unique, and
// it is strongly encouraged that the number be unique as well.
// Passing a -1 for the number will select the next valid number and
// assign it to this entity.
func (c *Client) EntityCreate(ctx context.Context, id, secret string, number int) error {
	if err := c.makeWritable(); err != nil {
		return err
	}

	ctx = c.appendMetadata(ctx)
	r := rpc.EntityRequest{
		Entity: &pb.Entity{
			ID:     &id,
			Secret: &secret,
			Number: proto.Int32(int32(number)),
		},
	}
	_, err := c.rpc.EntityCreate(ctx, &r)
	return err
}

// EntityUpdate alters the generic metadata on an existing entity.  It
// cannot modify keys or untyped metadata.
func (c *Client) EntityUpdate(ctx context.Context, id string, meta *pb.EntityMeta) error {
	if err := c.makeWritable(); err != nil {
		return err
	}

	ctx = c.appendMetadata(ctx)
	r := rpc.EntityRequest{
		Data: &pb.Entity{
			ID:   &id,
			Meta: meta,
		},
	}
	_, err := c.rpc.EntityUpdate(ctx, &r)
	return err
}

// EntityInfo returns information about an entity.  This function does
// not require authentication, and can be performed with an
// unauthenticated context.
func (c *Client) EntityInfo(ctx context.Context, id string) (pb.Entity, error) {
	ctx = c.appendMetadata(ctx)
	r := rpc.EntityRequest{
		Entity: &pb.Entity{
			ID: &id,
		},
	}

	res, err := c.rpc.EntityInfo(ctx, &r)
	if err != nil {
		return pb.Entity{}, err
	}
	return *res.GetEntities()[0], nil
}

// EntitySearch performs a search of all entities.  This search will
// return a slice of zero or more entities that matched the search
// criteria.  Searching does not require an authenticated context.
func (c *Client) EntitySearch(ctx context.Context, expr string) ([]*pb.Entity, error) {
	ctx = c.appendMetadata(ctx)
	r := rpc.SearchRequest{
		Expression: &expr,
	}

	res, err := c.rpc.EntitySearch(ctx, &r)
	return res.GetEntities(), err
}

// EntityUM handles operations concerning the untyped key-value store
// on each entity.  This data is not directly processed by NetAuth or
// visible in search indexes, but is useful for integrating with 3rd
// party systems as it provides an ideal place to store alternate keys
// or IDs.  Reads may be performed without authentication, writes must
// be authenticated.
func (c *Client) EntityUM(ctx context.Context, target, action, key, value string) (map[string][]string, error) {
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

	res, err := c.rpc.EntityUM(ctx, &r)
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

// EntityKVGet returns the values for a key if it exists.
func (c *Client) EntityKVGet(ctx context.Context, id, key string) (map[string][]string, error) {
	ctx = c.appendMetadata(ctx)
	r := rpc.KV2Request{
		Target: &id,
		Data: &pb.KVData{
			Key: &key,
		},
	}

	res, err := c.rpc.EntityKVGet(ctx, &r)
	if err != nil {
		return nil, err
	}

	out := make(map[string][]string, len(res.GetKVData()))
	for _, kvd := range res.GetKVData() {
		sort.Slice(kvd.Values, func(i, j int) bool {
			return kvd.Values[i].GetIndex() < kvd.Values[j].GetIndex()
		})
		for _, v := range kvd.GetValues() {
			out[kvd.GetKey()] = append(out[kvd.GetKey()], v.GetValue())
		}
	}

	return out, nil
}

// EntityKVAdd adds a single key to the specified entity.  The key
// specified must not already exist.  The order values are provided
// will be preserved.
func (c *Client) EntityKVAdd(ctx context.Context, id, key string, values []string) error {
	if err := c.makeWritable(); err != nil {
		return err
	}
	ctx = c.appendMetadata(ctx)

	r := rpc.KV2Request{
		Target: &id,
		Data: &pb.KVData{
			Key: &key,
		},
	}

	v := make([]*pb.KVValue, len(values))
	for i := range values {
		v[i] = &pb.KVValue{
			Value: &values[i],
			Index: proto.Int32(int32(i)),
		}
	}
	r.Data.Values = v

	_, err := c.rpc.EntityKVAdd(ctx, &r)
	return err
}

// EntityKVDel deletes a single existing key from the target.
func (c *Client) EntityKVDel(ctx context.Context, id, key string) error {
	if err := c.makeWritable(); err != nil {
		return err
	}
	ctx = c.appendMetadata(ctx)

	r := rpc.KV2Request{
		Target: &id,
		Data: &pb.KVData{
			Key: &key,
		},
	}
	_, err := c.rpc.EntityKVDel(ctx, &r)
	return err
}

// EntityKVReplace replaces a the values for a single key that must
// already exist.  Similar to add, ordering will be preserved.
func (c *Client) EntityKVReplace(ctx context.Context, id, key string, values []string) error {
	if err := c.makeWritable(); err != nil {
		return err
	}
	ctx = c.appendMetadata(ctx)

	r := rpc.KV2Request{
		Target: &id,
		Data: &pb.KVData{
			Key: &key,
		},
	}

	v := make([]*pb.KVValue, len(values))
	for i := range values {
		v[i] = &pb.KVValue{
			Value: &values[i],
			Index: proto.Int32(int32(i)),
		}
	}
	r.Data.Values = v

	_, err := c.rpc.EntityKVReplace(ctx, &r)
	return err
}

// EntityKeys handles updates to public keys stored on an entity.
// These keys are public and can be queried without authentication.
// The idea is to provide a means of distributing public keys for SSH
// and PGP.
func (c *Client) EntityKeys(ctx context.Context, id, action, ktype, key string) (map[string][]string, error) {
	if strings.ToUpper(action) != "READ" {
		if err := c.makeWritable(); err != nil {
			return nil, err
		}
	}

	ctx = c.appendMetadata(ctx)
	action = strings.ToUpper(action)
	a, ok := rpc.Action_value[action]
	if !ok {
		return nil, errors.New("action must be one of ADD, DROP, READ")
	}

	r := rpc.KVRequest{
		Target: &id,
		Action: rpc.Action(a).Enum(),
		Key:    &ktype,
		Value:  &key,
	}

	res, err := c.rpc.EntityKeys(ctx, &r)
	if err != nil {
		return nil, err
	}

	if action == "READ" {
		kv := parseKV(res.GetStrings())
		if ktype != "*" {
			// Asked for a specific key, fish it out and
			// return a much sparser map.
			return map[string][]string{ktype: kv[ktype]}, nil
		}
		return kv, nil
	}
	return nil, nil
}

// EntityDestroy is used to permanently remove entities from the
// server.  This is not recommended and should not be done without
// good reason.  The best practice is to instead have a group that
// defunct entities get moved to and then locked.  This will prevent
// authentication, while maintaining integrity of the backing tree.
// This function does not maintain referential integrity, so be
// careful about removing the last standing admin of a particular
// type.
func (c *Client) EntityDestroy(ctx context.Context, id string) error {
	if err := c.makeWritable(); err != nil {
		return err
	}

	ctx = c.appendMetadata(ctx)
	r := rpc.EntityRequest{
		Entity: &pb.Entity{
			ID: &id,
		},
	}

	_, err := c.rpc.EntityDestroy(ctx, &r)
	return err
}

// EntityLock sets the lock bit on the provided entity which will
// effectively prevent authentication from proceeding even if correct
// authentication information is provided.
func (c *Client) EntityLock(ctx context.Context, id string) error {
	if err := c.makeWritable(); err != nil {
		return err
	}

	ctx = c.appendMetadata(ctx)
	r := rpc.EntityRequest{
		Entity: &pb.Entity{
			ID: &id,
		},
	}
	_, err := c.rpc.EntityLock(ctx, &r)
	return err
}

// EntityUnlock is the inverse of EntityLock.  See EntityLock for more
// information.
func (c *Client) EntityUnlock(ctx context.Context, id string) error {
	if err := c.makeWritable(); err != nil {
		return err
	}

	ctx = c.appendMetadata(ctx)
	r := rpc.EntityRequest{
		Entity: &pb.Entity{
			ID: &id,
		},
	}
	_, err := c.rpc.EntityUnlock(ctx, &r)
	return err
}

// EntityGroups returns the effective group membership of the named entity.
func (c *Client) EntityGroups(ctx context.Context, id string) ([]*pb.Group, error) {
	ctx = c.appendMetadata(ctx)
	r := rpc.EntityRequest{
		Entity: &pb.Entity{
			ID: &id,
		},
	}
	res, err := c.rpc.EntityGroups(ctx, &r)
	return res.GetGroups(), err
}
