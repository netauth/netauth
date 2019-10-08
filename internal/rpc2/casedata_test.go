package rpc2

import (
	"context"

	"google.golang.org/grpc/metadata"

	"github.com/NetAuth/NetAuth/internal/token/null"

	pb "github.com/NetAuth/Protocol/v2"
)

var (
	ValidAuthData = &pb.AuthData{
		Token: &null.ValidToken,
	}

	InvalidAuthData = &pb.AuthData{
		Token: &null.InvalidToken,
	}

	EmptyAuthData = &pb.AuthData{
		Token: &null.ValidEmptyToken,
	}

	PrivilegedContext      = metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", null.ValidToken))
	UnprivilegedContext    = metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", null.ValidEmptyToken))
	UnauthenticatedContext = metadata.NewIncomingContext(context.Background(), nil)
	InvalidAuthContext     = metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", null.InvalidToken))
)
