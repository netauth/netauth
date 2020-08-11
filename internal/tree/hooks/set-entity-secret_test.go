package hooks

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-hclog"

	"github.com/netauth/netauth/internal/crypto/nocrypto"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

func TestSetEntitySecret(t *testing.T) {
	crypt, err := nocrypto.New(hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}

	hook, err := NewSetEntitySecret(tree.RefContext{Crypto: crypt})
	if err != nil {
		t.Fatal(err)
	}

	e := &pb.Entity{}
	de := &pb.Entity{Secret: proto.String("security")}

	if err := hook.Run(e, de); err != nil {
		t.Fatal(err)
	}

	if e.GetSecret() != "security" {
		t.Log(e)
		t.Fatal("Spec error - please trace hook")
	}
}

func TestSetEntitySecretCB(t *testing.T) {
	setEntitySecretCB()
}
