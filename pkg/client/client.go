package client

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/grpc"

	pb "github.com/NetAuth/NetAuth/proto"
)

func NewClient(server string, port int) (pb.NetAuthClient, error) {
	// Setup the connection and defer the close.
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", server, port), grpc.WithInsecure())

	// Create a client to use later on.
	return pb.NewNetAuthClient(conn), err
}

func Ping(server string, port int, clientID string) (string, error) {
	request := new(pb.PingRequest)
	request.ClientID = ensureClientID(clientID)

	client, err := NewClient(server, port)
	if err != nil {
		return "", err
	}
	pingResult, err := client.Ping(context.Background(), request)
	return pingResult.GetMsg(), nil
}

func Authenticate(server string, port int, clientID string, serviceID string, entity string, secret string) (string, error) {
	e := new(pb.Entity)
	e.ID = &entity
	e.Secret = &secret

	request := new(pb.NetAuthRequest)
	request.Entity = e
	request.ClientID = ensureClientID(clientID)
	request.ServiceID = ensureServiceID(serviceID)

	c, err := NewClient(server, port)
	if err != nil {
		return "", err
	}
	authResult, err := c.AuthEntity(context.Background(), request)
	if err != nil {
		return "", err
	}

	return authResult.GetMsg(), nil
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
