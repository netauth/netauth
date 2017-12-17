package rpc

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/golang/protobuf/proto"

	pb "github.com/NetAuth/NetAuth/proto"
)

func (s *NetAuthServer) Ping(ctx context.Context, pingRequest *pb.PingRequest) (*pb.PingResponse, error) {
	// Ping takes in a request from the client, and then replies
	// with a Pong containing the server status.

	log.Printf("Ping from %s", pingRequest.GetClientID())

	reply := new(pb.PingResponse)
	hostname, err := os.Hostname()
	if err != nil {
		log.Println("Hostname could not be determined!")
		hostname = "BOGUS_HOST"
	}
	reply.Msg = proto.String(fmt.Sprintf("NetAuth server on %s is ready to serve!", hostname))
	return reply, nil
}
