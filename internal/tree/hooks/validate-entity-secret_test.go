package hooks

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-hclog"

	"github.com/netauth/netauth/internal/crypto/nocrypto"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

func TestValidateEntitySecret(t *testing.T) {
	crypt, err := nocrypto.New(hclog.NewNullLogger())
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
