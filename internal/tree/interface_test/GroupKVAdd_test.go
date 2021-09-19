package interface_test

import (
	"testing"

	"google.golang.org/protobuf/proto"

	pb "github.com/netauth/protocol"
)

func TestGroupKVAdd(t *testing.T) {
	m, ctx := newTreeManager(t)

	addGroup(t, ctx)

	kv1 := []*pb.KVData{{
		Key: proto.String("key1"),
		Values: []*pb.KVValue{{
			Value: proto.String("value1"),
		}},
	}}

	if err := m.GroupKVAdd("group1", kv1); err != nil {
		t.Fatal(err)
	}

	kvtest, err := m.GroupKVGet("group1", kv1)
	if err != nil {
		t.Fatal(err)
	}

	if !proto.Equal(kvtest[0], kv1[0]) {
		t.Error("Set a key and got different data back")
	}
}
