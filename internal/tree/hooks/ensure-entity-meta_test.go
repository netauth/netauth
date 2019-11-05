package hooks

import (
	"testing"

	"github.com/netauth/netauth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

func TestEnsureEntityMeta(t *testing.T) {
	hook, err := NewEnsureEntityMeta(tree.RefContext{})
	if err != nil {
		t.Fatal(err)
	}

	e := &pb.Entity{}
	if err := hook.Run(e, &pb.Entity{}); err != nil {
		t.Fatal(err)
	}

	if e.Meta == nil {
		t.Log(e)
		t.Fatal("Spec error - please trace hook")
	}
}
