package client

import (
	"fmt"
	"log"

	"google.golang.org/grpc"

	pb "github.com/NetAuth/NetAuth/proto"
)

func NewClient(server string, port int) pb.NetAuthClient {
	// Setup the connection and defer the close.
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", server, port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to NetAuth: %s", err)
	}
	// This success message is very misleading, in theory if you
	// see this then a valid connection has been made to the
	// server, but this isn't really the case.  This will show
	// connected in all cases where the Dial() function returns
	// successfully, whether or not it has actually connected to
	// the NetAuth service is another matter entirely.  In theory
	// we could fire a PingRequest() before printing this message,
	// but that's somewhat superfluous when this will all fail out
	// within the next second if there's a problem.
	log.Printf("Connected to NetAuth server at %s:%d", server, port)

	// Create a client to use later on.
	return pb.NewNetAuthClient(conn)
}
