package hooks

import (
	"context"
	"testing"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/protobuf/proto"

	"github.com/netauth/netauth/internal/crypto/nocrypto"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

func TestValidateEntitySecret(t *testing.T) {
	crypt, err := nocrypto.New(hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}

	hook, err := NewValidateEntitySecret(tree.WithHookCrypto(crypt))
	if err != nil {
		t.Fatal(err)
	}

	e := &pb.Entity{Secret: proto.String("secret")}
	de := &pb.Entity{Secret: proto.String("secret")}
	if err := hook.Run(context.Background(), e, de); err != nil {
		t.Fatal(err)
	}
}

func TestValidateEntitySecretCB(t *testing.T) {
	validateEntitySecretCB()
}
