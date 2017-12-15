package client

import (
	"context"
	"fmt"
	"log"
	"os"

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

func Ping(server string, port int, clientID string) bool {
	log.Printf("Pinging the server")

	request := new(pb.PingRequest)
	request.ClientID = ensureClientID(clientID)

	client := NewClient(server, port)
	pingResult, err := client.Ping(context.Background(), request)
	if err != nil {
		log.Fatalf("Ping failed: %s", err)
	}
	log.Printf("%s", pingResult)
	return pingResult.GetHealthy()
}

func Authenticate(server string, port int, clientID string, serviceID string, entity string, secret string) bool {
	log.Printf("Trying to authenticate %s", entity)

	e := new(pb.Entity)
	e.ID = &entity
	e.Secret = &secret

	request := new(pb.NetAuthRequest)
	request.Entity = e
	request.ClientID = ensureClientID(clientID)
	request.ServiceID = ensureServiceID(serviceID)

	c := NewClient(server, port)
	authResult, err := c.AuthEntity(context.Background(), request)
	if err != nil {
		log.Fatalf("Could not auth: %s", err)
	}
	if authResult == nil {
		log.Fatal("recieved nil reply for AuthEntity()")
	}
	log.Printf("%v", authResult)
	return authResult.GetSuccess()
}

func ensureClientID(clientID string) *string {
	if clientID == "" {
		hostname, err := os.Hostname()
		if err != nil {
			clientID = "BOGUS_CLIENT"
			return &clientID
		}
		clientID = hostname
	}
	return &clientID
}

func ensureServiceID(serviceID string) *string {
	if serviceID == "" {
		serviceID = "BOGUS_SERVICE"
	}
	return &serviceID
}
