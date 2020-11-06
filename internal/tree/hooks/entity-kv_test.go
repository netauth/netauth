package hooks

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

func TestEntityKVAdd(t *testing.T) {
	entitykvcb()

	cases := []struct {
		e       *pb.Entity
		de      *pb.Entity
		wantErr error
	}{
		{
			e: &pb.Entity{
				Meta: &pb.EntityMeta{
					KV: []*pb.KVData{
						{
							Key: proto.String("key1"),
							Values: []*pb.KVValue{
								{Value: proto.String("value1")},
							},
						},
					},
				},
			},
			de: &pb.Entity{
				Meta: &pb.EntityMeta{
					KV: []*pb.KVData{
						{
							Key: proto.String("key1"),
							Values: []*pb.KVValue{
								{Value: proto.String("value1")},
							},
						},
					},
				},
			},
			wantErr: tree.ErrKeyExists,
		},
		{
			e: &pb.Entity{
				Meta: &pb.EntityMeta{
					KV: []*pb.KVData{
						{
							Key: proto.String("key1"),
							Values: []*pb.KVValue{
								{Value: proto.String("value1")},
							},
						},
					},
				},
			},
			de: &pb.Entity{
				Meta: &pb.EntityMeta{
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
			},
			wantErr: nil,
		},
		{
			e: &pb.Entity{
				Meta: &pb.EntityMeta{
					KV: []*pb.KVData{
						{
							Key: proto.String("key1"),
							Values: []*pb.KVValue{
								{Value: proto.String("value1")},
							},
						},
					},
				},
			},
			de:      &pb.Entity{},
			wantErr: tree.ErrFailedPrecondition,
		},
	}

	h, _ := newEntityKVAdd(tree.RefContext{})

	for i, c := range cases {
		if err := h.Run(c.e, c.de); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}

func TestEntityKVDel(t *testing.T) {
	entitykvcb()

	cases := []struct {
		e       *pb.Entity
		de      *pb.Entity
		wantErr error
	}{
		{
			e: &pb.Entity{
				Meta: &pb.EntityMeta{
					KV: []*pb.KVData{
						{
							Key: proto.String("key1"),
							Values: []*pb.KVValue{
								{Value: proto.String("value1")},
							},
						},
					},
				},
			},
			de: &pb.Entity{
				Meta: &pb.EntityMeta{
					KV: []*pb.KVData{
						{
							Key: proto.String("key1"),
							Values: []*pb.KVValue{
								{Value: proto.String("value1")},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			e: &pb.Entity{
				Meta: &pb.EntityMeta{
					KV: []*pb.KVData{
						{
							Key: proto.String("key1"),
							Values: []*pb.KVValue{
								{Value: proto.String("value1")},
							},
						},
					},
				},
			},
			de: &pb.Entity{
				Meta: &pb.EntityMeta{
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
			},
			wantErr: tree.ErrNoSuchKey,
		},
		{
			e: &pb.Entity{
				Meta: &pb.EntityMeta{
					KV: []*pb.KVData{
						{
							Key: proto.String("key1"),
							Values: []*pb.KVValue{
								{Value: proto.String("value1")},
							},
						},
					},
				},
			},
			de:      &pb.Entity{},
			wantErr: tree.ErrFailedPrecondition,
		},
	}

	h, _ := newEntityKVDel(tree.RefContext{})

	for i, c := range cases {
		if err := h.Run(c.e, c.de); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}

func TestEntityKVReplace(t *testing.T) {
	entitykvcb()

	cases := []struct {
		e       *pb.Entity
		de      *pb.Entity
		wantErr error
	}{
		{
			e: &pb.Entity{
				Meta: &pb.EntityMeta{
					KV: []*pb.KVData{
						{
							Key: proto.String("key1"),
							Values: []*pb.KVValue{
								{Value: proto.String("value1")},
							},
						},
					},
				},
			},
			de: &pb.Entity{
				Meta: &pb.EntityMeta{
					KV: []*pb.KVData{
						{
							Key: proto.String("key1"),
							Values: []*pb.KVValue{
								{Value: proto.String("value1")},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			e: &pb.Entity{
				Meta: &pb.EntityMeta{
					KV: []*pb.KVData{
						{
							Key: proto.String("key1"),
							Values: []*pb.KVValue{
								{Value: proto.String("value1")},
							},
						},
					},
				},
			},
			de: &pb.Entity{
				Meta: &pb.EntityMeta{
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
			},
			wantErr: tree.ErrNoSuchKey,
		},
	}

	h, _ := newEntityKVReplace(tree.RefContext{})

	for i, c := range cases {
		if err := h.Run(c.e, c.de); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}
