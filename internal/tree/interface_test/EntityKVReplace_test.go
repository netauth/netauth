package interface_test

import (
	"testing"

	"github.com/golang/protobuf/proto"

	pb "github.com/netauth/protocol"
)

func TestEntityKVReplace(t *testing.T) {
	m, ctx := newTreeManager(t)

	addEntity(t, ctx)

	kv1 := []*pb.KVData{{
		Key: proto.String("key1"),
		Values: []*pb.KVValue{{
			Value: proto.String("value1"),
		}},
	}}
	kv2 := []*pb.KVData{{
		Key: proto.String("key1"),
		Values: []*pb.KVValue{
			{Value: proto.String("value1")},
			{Value: proto.String("value2")},
		},
	}}

	if err := m.EntityKVAdd("entity1", kv1); err != nil {
		t.Fatal(err)
	}

	kvtest, err := m.EntityKVGet("entity1", kv1)
	if err != nil {
		t.Fatal(err)
	}

	if !proto.Equal(kvtest[0], kv1[0]) {
		t.Error("Set a key and got different data back")
	}

	if err := m.EntityKVReplace("entity1", kv2); err != nil {
		t.Fatal(err)
	}

	// We can do the get on kv1 here because the name should be
	// the same for both, and so we do a get on kv1 and expect the
	// results to be kv2.
	res, err := m.EntityKVGet("entity1", kv1)
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 || len(res[0].GetValues()) != 2 {
		t.Log(res)
		t.Fatal("key was not updated")
	}
}
