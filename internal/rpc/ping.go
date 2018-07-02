package rpc

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/NetAuth/NetAuth/internal/health"

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
	log.Printf("Ping from %s@%s", client.GetService(), client.GetID())

	reply := new(pb.PingResponse)
	reply.Healthy = proto.Bool(health.Get())
	hostname, err := os.Hostname()
	if err != nil {
		log.Println("Hostname could not be determined!")
		hostname = "BOGUS_HOST"
	}

	if *reply.Healthy {
		reply.Msg = proto.String(fmt.Sprintf("NetAuth server on %s is ready to serve!", hostname))
	} else {
		reply.Msg = proto.String(fmt.Sprintf("NetAuth server on %s is not ready to serve at this time!", hostname))
	}
	return reply, toWireError(nil)
}
