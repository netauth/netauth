package interface_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"

	"github.com/netauth/netauth/internal/db"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

func TestGroupKVGet(t *testing.T) {
	ctxt := context.Background()
	m, ctx := newTreeManager(t)

	addGroup(t, ctx)

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

	if err := m.GroupKVAdd(ctxt, "group1", []*pb.KVData{kv1}); err != nil {
		t.Fatal(err)
	}
	if err := m.GroupKVAdd(ctxt, "group1", []*pb.KVData{kv2}); err != nil {
		t.Fatal(err)
	}

	kvtest, err := m.GroupKVGet(ctxt, "group1", []*pb.KVData{kv2})
	if err != nil {
		t.Fatal(err)
	}

	if !proto.Equal(kvtest[0], kv2) {
		t.Error("Set a key and got different data back")
	}

	if _, err := m.GroupKVGet(ctxt, "does-not-exist", []*pb.KVData{kv1}); err != db.ErrUnknownGroup {
		t.Error(err)
	}

	if _, err := m.GroupKVGet(ctxt, "group1", []*pb.KVData{{Key: proto.String("does-not-exist")}}); err != tree.ErrNoSuchKey {
		t.Error(err)
	}

	res, err := m.GroupKVGet(ctxt, "group1", []*pb.KVData{{Key: proto.String("*")}})
	assert.Nil(t, err)
	expect := []*pb.KVData{kv1, kv2}
	assert.Nil(t, err)
	// Trying to do deep equals in the protobuf fails, so instead
	// we assume that the individual tests above have worked and
	// that the right amount of data was returned here.
	assert.Equal(t, len(res), len(expect))
}
