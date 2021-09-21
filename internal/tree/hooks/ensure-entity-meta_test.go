package hooks

import (
	"context"
	"testing"

	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

func TestEnsureEntityMeta(t *testing.T) {
	hook, err := NewEnsureEntityMeta(tree.RefContext{})
	if err != nil {
		t.Fatal(err)
	}

	e := &pb.Entity{}
	if err := hook.Run(context.Background(), e, &pb.Entity{}); err != nil {
		t.Fatal(err)
	}

	if e.Meta == nil {
		t.Log(e)
		t.Fatal("Spec error - please trace hook")
	}
}

func TestEnsureEntityMetaCB(t *testing.T) {
	ensureEntityMetaCB()
}
