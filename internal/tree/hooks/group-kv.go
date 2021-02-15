package hooks

import (
	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

func init() {
	startup.RegisterCallback(groupkvcb)
}

// GroupKV handles additions of a new key to the group structure.
type GroupKV struct {
	tree.BaseHook

	do func(*pb.Group, *pb.Group) error
}

// Run proxies to the do function which is set based on what the hook
// is supposed to do.
func (ekv *GroupKV) Run(g, dg *pb.Group) error {
	return ekv.do(g, dg)
}

// add will iterate through all keys and check if the key already
// exists.  If it does not it will be appended to the end of the list
// of keys.
func (ekv *GroupKV) add(g, dg *pb.Group) error {
	if len(dg.GetKV()) != 1 {
		return tree.ErrFailedPrecondition
	}
	compare := dg.GetKV()[0].GetKey()

	for _, k := range g.GetKV() {
		if k.GetKey() == compare {
			return tree.ErrKeyExists
		}
	}
	g.KV = append(g.KV, dg.GetKV()[0])
	return nil
}

// del looks for the specified key and tries to remove it, returning
// an error if it does not exist.
func (ekv *GroupKV) del(g, dg *pb.Group) error {
	if len(dg.GetKV()) != 1 {
		return tree.ErrFailedPrecondition
	}
	compare := dg.GetKV()[0].GetKey()

	out := []*pb.KVData{}
	for _, k := range g.GetKV() {
		if k.GetKey() == compare {
			continue
		}
		out = append(out, k)
	}

	if len(out) == len(g.GetKV()) {
		return tree.ErrNoSuchKey
	}
	g.KV = out
	return nil
}

// replace looks for an existing key, and then replaces it.  It is a
// convenience function on top of a delete/add paired call.
func (ekv *GroupKV) replace(g, dg *pb.Group) error {
	if err := ekv.del(g, dg); err != nil {
		return err
	}
	return ekv.add(g, dg)
}

func newGroupKVAdd(c tree.RefContext) (tree.GroupHook, error) {
	x := &GroupKV{}
	x.BaseHook = tree.NewBaseHook("kv-add", 50)
	x.do = x.add
	return x, nil
}

func newGroupKVDel(c tree.RefContext) (tree.GroupHook, error) {
	x := &GroupKV{}
	x.BaseHook = tree.NewBaseHook("kv-del", 50)
	x.do = x.del
	return x, nil
}

func newGroupKVReplace(c tree.RefContext) (tree.GroupHook, error) {
	x := &GroupKV{}
	x.BaseHook = tree.NewBaseHook("kv-replace", 50)
	x.do = x.replace
	return x, nil
}

func groupkvcb() {
	tree.RegisterGroupHookConstructor("kv-add", newGroupKVAdd)
	tree.RegisterGroupHookConstructor("kv-del", newGroupKVDel)
	tree.RegisterGroupHookConstructor("kv-replace", newGroupKVReplace)
}
