package rpc

import (
	"context"
	"fmt"
	"os"

	"github.com/netauth/netauth/internal/health"

	"github.com/golang/protobuf/proto"

	pb "github.com/NetAuth/Protocol"
)

// Ping requests the health status of the server and returns it to the
// client.  This is designed to be a virtually free action that should
// be safe to invoke at any time to see if the server is available.
func (s *NetAuthServer) Ping(ctx context.Context, pingRequest *pb.PingRequest) (*pb.PingResponse, error) {
	// Ping takes in a request from the client, and then replies
	// with a Pong containing the server status.

	client := pingRequest.GetInfo()
	s.Log.Info("Ping", "service", client.GetService(), "client", client.GetID())

	reply := new(pb.PingResponse)
	status := health.Check()
	reply.Healthy = proto.Bool(status.OK)
	hostname, err := os.Hostname()
	if err != nil {
		s.Log.Warn("Hostname could not be determined!")
		hostname = "BOGUS_HOST"
	}

	st := "ready"
	if !status.OK {
		st = "not ready"
	}
	msg := fmt.Sprintf("NetAuth server on %s is %s\n\n%s", hostname, st, status)

	reply.Msg = &msg
	return reply, toWireError(nil)
}
