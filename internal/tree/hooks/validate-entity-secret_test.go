package hooks

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/internal/crypto/nocrypto"
	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

func TestValidateEntitySecret(t *testing.T) {
	crypt, err := nocrypto.New()
	if err != nil {
		t.Fatal(err)
	}

	hook, err := NewValidateEntitySecret(tree.RefContext{Crypto: crypt})
	if err != nil {
		t.Fatal(err)
	}

	e := &pb.Entity{Secret: proto.String("secret")}
	de := &pb.Entity{Secret: proto.String("secret")}
	if err := hook.Run(e, de); err != nil {
		t.Fatal(err)
	}
}
