package interface_test

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"

	"github.com/netauth/netauth/internal/db"

	pb "github.com/netauth/protocol"
)

func TestEntityKVGet(t *testing.T) {
	m, ctx := newTreeManager(t)

	addEntity(t, ctx)

	kv1 := &pb.KVData{
		Key: proto.String("key1"),
		Values: []*pb.KVValue{
			{Value: proto.String("value1")},
		},
	}
	kv2 := &pb.KVData{
		Key: proto.String("key2"),
		Values: []*pb.KVValue{
			{Value: proto.String("value2")},
			{Value: proto.String("value15")},
			{Value: proto.String("value09")},
		},
	}

	if err := m.EntityKVAdd("entity1", []*pb.KVData{kv1}); err != nil {
		t.Fatal(err)
	}
	if err := m.EntityKVAdd("entity1", []*pb.KVData{kv2}); err != nil {
		t.Fatal(err)
	}

	kvtest, err := m.EntityKVGet("entity1", []*pb.KVData{kv2})
	if err != nil {
		t.Fatal(err)
	}

	if !proto.Equal(kvtest[0], kv2) {
		t.Error("Set a key and got different data back")
	}

	if _, err := m.EntityKVGet("does-not-exist", []*pb.KVData{kv1}); err != db.ErrUnknownEntity {
		t.Error(err)
	}

	res, err := m.EntityKVGet("entity1", []*pb.KVData{{Key: proto.String("*")}})
	assert.Nil(t, err)
	assert.Equal(t, res, []*pb.KVData{kv1, kv2})
}
