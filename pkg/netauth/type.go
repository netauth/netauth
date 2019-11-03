package netauth

import (
	"github.com/hashicorp/go-hclog"

	"github.com/NetAuth/NetAuth/internal/token"

	rpc "github.com/NetAuth/Protocol/v2"
)

// Client is an RPC client shim that makes communicating with the
// NetAuth server easier.  The client has helpers for attaching
// parameters to the request, for crafting protobufs, and for handling
// other common tasks.
type Client struct {
	TokenCache
	token.Service

	rpc rpc.NetAuth2Client
	log hclog.Logger

	clientName  string
	serviceName string

	writeable bool
}
