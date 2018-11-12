package hooks

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

func TestValidateEntityUnlocked(t *testing.T) {
	hook, err := NewValidateEntityUnlocked(tree.RefContext{})
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		e       *pb.Entity
		wantErr error
	}{
		{&pb.Entity{Meta: &pb.EntityMeta{Locked: proto.Bool(true)}}, tree.ErrEntityLocked},
		{&pb.Entity{Meta: &pb.EntityMeta{Locked: proto.Bool(false)}}, nil},
	}

	for i, c := range cases {
		if err := hook.Run(c.e, &pb.Entity{}); err != c.wantErr {
			t.Errorf("Case %d - Got: %v Want: %v", i, err, c.wantErr)
		}
	}
}
