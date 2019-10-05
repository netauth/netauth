package rpc2

import (
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
)
