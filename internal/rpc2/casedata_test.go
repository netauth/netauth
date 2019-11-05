package rpc2

import (
	"context"

	"google.golang.org/grpc/metadata"

	"github.com/netauth/netauth/internal/token/null"
)

var (
	PrivilegedContext      = metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", null.ValidToken))
	UnprivilegedContext    = metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", null.ValidEmptyToken))
	UnauthenticatedContext = metadata.NewIncomingContext(context.Background(), nil)
	InvalidAuthContext     = metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", null.InvalidToken))
)
