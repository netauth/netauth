package hooks

import (
	"context"

	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

func init() {
	startup.RegisterCallback(entitykvcb)
}

// EntityKV handles additions of a new key to the entity structure.
type EntityKV struct {
	tree.BaseHook

	do func(*pb.Entity, *pb.Entity) error
}

// Run proxies to the do function which is set based on what the hook
// is supposed to do.
func (ekv *EntityKV) Run(_ context.Context, e, de *pb.Entity) error {
	return ekv.do(e, de)
}

// add will iterate through all keys and check if the key already
// exists.  If it does not it will be appended to the end of the list
// of keys.
func (ekv *EntityKV) add(e, de *pb.Entity) error {
	if len(de.GetMeta().GetKV()) != 1 {
		return tree.ErrFailedPrecondition
	}
	compare := de.GetMeta().GetKV()[0].GetKey()

	for _, k := range e.GetMeta().GetKV() {
		if k.GetKey() == compare {
			return tree.ErrKeyExists
		}
	}
	e.Meta.KV = append(e.Meta.KV, de.GetMeta().GetKV()[0])
	return nil
}

// del looks for the specified key and tries to remove it, returning
// an error if it does not exist.
func (ekv *EntityKV) del(e, de *pb.Entity) error {
	if len(de.GetMeta().GetKV()) != 1 {
		return tree.ErrFailedPrecondition
	}
	compare := de.GetMeta().GetKV()[0].GetKey()

	out := []*pb.KVData{}
	for _, k := range e.GetMeta().GetKV() {
		if k.GetKey() == compare {
			continue
		}
		out = append(out, k)
	}

	if len(out) == len(e.GetMeta().GetKV()) {
		return tree.ErrNoSuchKey
	}
	e.Meta.KV = out
	return nil
}

// replace looks for an existing key, and then replaces it.  It is a
// convenience function on top of a delete/add paired call.
func (ekv *EntityKV) replace(e, de *pb.Entity) error {
	if err := ekv.del(e, de); err != nil {
		return err
	}
	return ekv.add(e, de)
}

func newEntityKVAdd(opts ...tree.HookOption) (tree.EntityHook, error) {
	opts = append([]tree.HookOption{
		tree.WithHookName("kv-add"),
		tree.WithHookPriority(50),
	}, opts...)
	x := &EntityKV{}
	x.BaseHook = tree.NewBaseHook(opts...)
	x.do = x.add
	return x, nil
}

func newEntityKVDel(opts ...tree.HookOption) (tree.EntityHook, error) {
	opts = append([]tree.HookOption{
		tree.WithHookName("kv-del"),
		tree.WithHookPriority(50),
	}, opts...)
	x := &EntityKV{}
	x.BaseHook = tree.NewBaseHook(opts...)
	x.do = x.del
	return x, nil
}

func newEntityKVReplace(opts ...tree.HookOption) (tree.EntityHook, error) {
	opts = append([]tree.HookOption{
		tree.WithHookName("kv-replace"),
		tree.WithHookPriority(50),
	}, opts...)
	x := &EntityKV{}
	x.BaseHook = tree.NewBaseHook(opts...)
	x.do = x.replace
	return x, nil
}

func entitykvcb() {
	tree.RegisterEntityHookConstructor("kv-add", newEntityKVAdd)
	tree.RegisterEntityHookConstructor("kv-del", newEntityKVDel)
	tree.RegisterEntityHookConstructor("kv-replace", newEntityKVReplace)
}
