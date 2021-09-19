package hooks

import (
	"testing"

	"github.com/netauth/netauth/internal/tree"
	"google.golang.org/protobuf/proto"

	pb "github.com/netauth/protocol"
)

func TestGroupKVAdd(t *testing.T) {
	groupkvcb()

	cases := []struct {
		e       *pb.Group
		de      *pb.Group
		wantErr error
	}{
		{
			e: &pb.Group{
				KV: []*pb.KVData{
					{
						Key: proto.String("key1"),
						Values: []*pb.KVValue{
							{Value: proto.String("value1")},
						},
					},
				},
			},
			de: &pb.Group{
				KV: []*pb.KVData{
					{
						Key: proto.String("key1"),
						Values: []*pb.KVValue{
							{Value: proto.String("value1")},
						},
					},
				},
			},
			wantErr: tree.ErrKeyExists,
		},
		{
			e: &pb.Group{
				KV: []*pb.KVData{
					{
						Key: proto.String("key1"),
						Values: []*pb.KVValue{
							{Value: proto.String("value1")},
						},
					},
				},
			},
			de: &pb.Group{
				KV: []*pb.KVData{
					{
						Key: proto.String("key2"),
						Values: []*pb.KVValue{
							{Value: proto.String("value2")},
							{Value: proto.String("value3")},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			e: &pb.Group{
				KV: []*pb.KVData{
					{
						Key: proto.String("key1"),
						Values: []*pb.KVValue{
							{Value: proto.String("value1")},
						},
					},
				},
			},
			de:      &pb.Group{},
			wantErr: tree.ErrFailedPrecondition,
		},
	}

	h, _ := newGroupKVAdd(tree.RefContext{})

	for i, c := range cases {
		if err := h.Run(c.e, c.de); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}

func TestGroupKVDel(t *testing.T) {
	groupkvcb()

	cases := []struct {
		e       *pb.Group
		de      *pb.Group
		wantErr error
	}{
		{
			e: &pb.Group{
				KV: []*pb.KVData{
					{
						Key: proto.String("key1"),
						Values: []*pb.KVValue{
							{Value: proto.String("value1")},
						},
					},
				},
			},
			de: &pb.Group{
				KV: []*pb.KVData{
					{
						Key: proto.String("key1"),
						Values: []*pb.KVValue{
							{Value: proto.String("value1")},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			e: &pb.Group{
				KV: []*pb.KVData{
					{
						Key: proto.String("key1"),
						Values: []*pb.KVValue{
							{Value: proto.String("value1")},
						},
					},
				},
			},
			de: &pb.Group{
				KV: []*pb.KVData{
					{
						Key: proto.String("key2"),
						Values: []*pb.KVValue{
							{Value: proto.String("value2")},
							{Value: proto.String("value3")},
						},
					},
				},
			},
			wantErr: tree.ErrNoSuchKey,
		},
		{
			e: &pb.Group{
				KV: []*pb.KVData{
					{
						Key: proto.String("key1"),
						Values: []*pb.KVValue{
							{Value: proto.String("value1")},
						},
					},
				},
			},
			de:      &pb.Group{},
			wantErr: tree.ErrFailedPrecondition,
		},
	}

	h, _ := newGroupKVDel(tree.RefContext{})

	for i, c := range cases {
		if err := h.Run(c.e, c.de); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}

func TestGroupKVReplace(t *testing.T) {
	groupkvcb()

	cases := []struct {
		e       *pb.Group
		de      *pb.Group
		wantErr error
	}{
		{
			e: &pb.Group{
				KV: []*pb.KVData{
					{
						Key: proto.String("key1"),
						Values: []*pb.KVValue{
							{Value: proto.String("value1")},
						},
					},
				},
			},
			de: &pb.Group{
				KV: []*pb.KVData{
					{
						Key: proto.String("key1"),
						Values: []*pb.KVValue{
							{Value: proto.String("value1")},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			e: &pb.Group{
				KV: []*pb.KVData{
					{
						Key: proto.String("key1"),
						Values: []*pb.KVValue{
							{Value: proto.String("value1")},
						},
					},
				},
			},
			de: &pb.Group{
				KV: []*pb.KVData{
					{
						Key: proto.String("key2"),
						Values: []*pb.KVValue{
							{Value: proto.String("value2")},
							{Value: proto.String("value3")},
						},
					},
				},
			},
			wantErr: tree.ErrNoSuchKey,
		},
	}

	h, _ := newGroupKVReplace(tree.RefContext{})

	for i, c := range cases {
		if err := h.Run(c.e, c.de); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}
