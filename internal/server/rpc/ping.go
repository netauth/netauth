package rpc

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/NetAuth/NetAuth/internal/server/health"

	"github.com/golang/protobuf/proto"

	pb "github.com/NetAuth/NetAuth/proto"
)

func (s *NetAuthServer) Ping(ctx context.Context, pingRequest *pb.PingRequest) (*pb.PingResponse, error) {
	// Ping takes in a request from the client, and then replies
	// with a Pong containing the server status.

	log.Printf("Ping from %s", pingRequest.GetClientID())

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
	return reply, nil
}
