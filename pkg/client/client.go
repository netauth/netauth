package client

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/NetAuth/NetAuth/internal/token"
	_ "github.com/NetAuth/NetAuth/internal/token/impl"

	"google.golang.org/grpc"

	pb "github.com/NetAuth/NetAuth/pkg/proto"
)

type netAuthClient struct {
	c          pb.NetAuthClient
	serviceID  *string
	clientID   *string
	tokenStore TokenStore

	tokenService token.TokenService
}

// New takes in the values that set up a client and builds a
// client.netAuthClient struct on which all other methods are bound.
// This drastically simplifies the construction of other functions.
func New(server string, port int, serviceID string, clientID string) (*netAuthClient, error) {
	// Setup the connection.
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", server, port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	// Get a tokenstore
	t, err := getTokenStore()
	if err != nil {
		log.Fatal(err)
	}

	// Get a token service
	ts, err := token.New()
	if err != nil {
		log.Fatal(err)
	}

	// Create a client to use later on.
	client := netAuthClient{
		c:            pb.NewNetAuthClient(conn),
		serviceID:    ensureServiceID(serviceID),
		clientID:     ensureClientID(clientID),
		tokenStore:   t,
		tokenService: ts,
	}

	return &client, nil
}

// Ping very simply pings the server.  The reply will contain the
// health status of the server as a server that replies and a server
// that can serve are two very different things (data might be
// reloading during the request).
func (n *netAuthClient) Ping() (string, error) {
	request := new(pb.PingRequest)
	request.Info = &pb.ClientInfo{
		ID:      n.clientID,
		Service: n.serviceID,
	}

	pingResult, err := n.c.Ping(context.Background(), request)
	if err != nil {
		return "RPC Error", err
	}
	return pingResult.GetMsg(), nil
}

// Authenticate takes in an entity and a secret and tries to validate
// that the identity is legitimate by verifying the secret provided.
func (n *netAuthClient) Authenticate(entity string, secret string) (string, error) {
	request := pb.NetAuthRequest{
		Entity: &pb.Entity{
			ID:     &entity,
			Secret: &secret,
		},
		Info: &pb.ClientInfo{
			ID:      n.clientID,
			Service: n.serviceID,
		},
	}

	authResult, err := n.c.AuthEntity(context.Background(), &request)
	if err != nil {
		return "", err
	}

	return authResult.GetMsg(), nil
}

// GetToken is identical to Authenticate except on success it will
// return a token which can be used to authorize additional later
// requests.
func (n *netAuthClient) GetToken(entity, secret string) (string, error) {
	// See if we have a local copy first.
	t, err := n.getTokenFromStore(entity)
	if err == nil {
		if _, err := n.InspectToken(t); err == nil {
			return t, nil
		}
	}

	request := pb.NetAuthRequest{
		Entity: &pb.Entity{
			ID:     &entity,
			Secret: &secret,
		},
		Info: &pb.ClientInfo{
			ID:      n.clientID,
			Service: n.serviceID,
		},
	}
	tokenResult, err := n.c.GetToken(context.Background(), &request)
	if err != nil {
		return "", err
	}

	t = tokenResult.GetToken()
	if err := n.tokenStore.DestroyToken(entity); err != nil {
		return "", err
	}
	err = n.putTokenInStore(entity, t)
	return t, err
}

// ValidateToken sends the token to the server for validation.  This
// is effectively asking the server to authenticate the token and not
// do anything else.  Returns a comment from the server and an error.
func (n *netAuthClient) ValidateToken(entity string) (string, error) {
	t, err := n.getTokenFromStore(entity)
	if err != nil {
		return "", err
	}

	request := pb.NetAuthRequest{
		Entity: &pb.Entity{
			ID: &entity,
		},
		AuthToken: &t,
		Info: &pb.ClientInfo{
			ID:      n.clientID,
			Service: n.serviceID,
		},
	}

	result, err := n.c.ValidateToken(context.Background(), &request)
	if err != nil {
		return "", err
	}

	return result.GetMsg(), nil
}

// ChangeSecret crafts a modEntity request with the correct fields to
// change an entity secret either via self authentication or via token
// authentication which is held by an appropriate administrator.
func (n *netAuthClient) ChangeSecret(e, s, me, ms, t string) (string, error) {
	request := pb.ModEntityRequest{
		Entity: &pb.Entity{
			ID:     &e,
			Secret: &s,
		},
		ModEntity: &pb.Entity{
			ID:     &me,
			Secret: &ms,
		},
		AuthToken: &t,
		Info: &pb.ClientInfo{
			ID:      n.clientID,
			Service: n.serviceID,
		},
	}

	result, err := n.c.ChangeSecret(context.Background(), &request)
	if err != nil {
		return "", err
	}
	return result.GetMsg(), nil
}

// NewEntity crafts a modEntity request with the correct fields to
// create a new entity.
func (n *netAuthClient) NewEntity(id string, uidn int32, secret, t string) (string, error) {
	request := pb.ModEntityRequest{
		Entity: &pb.Entity{
			ID:        &id,
			UidNumber: &uidn,
			Secret:    &secret,
		},
		AuthToken: &t,
		Info: &pb.ClientInfo{
			ID:      n.clientID,
			Service: n.serviceID,
		},
	}

	result, err := n.c.NewEntity(context.Background(), &request)
	if err != nil {
		return "", err
	}
	return result.GetMsg(), nil
}

// RemoveEntity removes an entity by the given name.  Only the
// 'entity' field of the modEntityRequest is required.
func (n *netAuthClient) RemoveEntity(id, token string) (string, error) {
	request := pb.ModEntityRequest{
		Entity: &pb.Entity{
			ID: &id,
		},
		AuthToken: &token,
		Info: &pb.ClientInfo{
			ID:      n.clientID,
			Service: n.serviceID,
		},
	}

	result, err := n.c.RemoveEntity(context.Background(), &request)
	if err != nil {
		return "", err
	}
	return result.GetMsg(), nil
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
