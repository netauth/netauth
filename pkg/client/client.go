package client

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/grpc"

	pb "github.com/NetAuth/NetAuth/proto"
)

func newClient(server string, port int) (pb.NetAuthClient, error) {
	// Setup the connection and defer the close.
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", server, port), grpc.WithInsecure())

	// Create a client to use later on.
	return pb.NewNetAuthClient(conn), err
}

func Ping(server string, port int, clientID string) (string, error) {
	request := new(pb.PingRequest)
	request.ClientID = ensureClientID(clientID)

	client, err := newClient(server, port)
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

	c, err := newClient(server, port)
	if err != nil {
		return "", err
	}
	authResult, err := c.AuthEntity(context.Background(), request)
	if err != nil {
		return "", err
	}

	return authResult.GetMsg(), nil
}

// NewEntity makes a request to the server to add a new entity.  It
// requires an existing entity for authentication and authorization to
// add the new one, as well as parameters to populate the core fields
// on the new entity.  This function returns a string message from the
// server and an error describing whether or not the server was able
// to add the requested entity.
func NewEntity(server string, port int, clientID string, serviceID string, entity string, secret string, newEntity string, newUIDNumber int32, newSecret string) (string, error) {
	// e is the entity that is requesting this change.  This
	// entity must have the correct capabilities to actually add a
	// new entity to the system.
	e := new(pb.Entity)
	e.ID = &entity
	e.Secret = &secret

	// ne is the new entity.  These fields are the ones that must
	// be set at the time of creation for a new entity.
	ne := new(pb.Entity)
	ne.ID = &newEntity
	ne.UidNumber = &newUIDNumber
	ne.Secret = &newSecret

	request := new(pb.ModEntityRequest)
	request.Entity = e
	request.ModEntity = ne

	c, err := newClient(server, port)
	if err != nil {
		return "", err
	}

	result, err := c.NewEntity(context.Background(), request)
	if err != nil {
		return "", err
	}

	return result.GetMsg(), nil
}

// RemoveEntity makes a request to the server to remove the named
// entity.  This must be authorized by an entity which has the
// appropriate capabilities to fulfill the remove request
func RemoveEntity(server string, port int, clientID string, serviceID string, entity string, secret string, delEntity string) (string, error) {
	// e is the entity requesting this change, it must have the
	// correct permissions to run the remove.
	e := new(pb.Entity)
	e.ID = &entity
	e.Secret = &secret

	// re is the entity to be removed
	re := new(pb.Entity)
	re.ID = &delEntity

	request := new(pb.ModEntityRequest)
	request.Entity = e
	request.ModEntity = re

	c, err := newClient(server, port)
	if err != nil {
		return "", err
	}

	result, err := c.RemoveEntity(context.Background(), request)
	if err != nil {
		return "", err
	}

	return result.GetMsg(), nil
}

// ChangeSecret changes the secret on the given modEntity under the
// authority of the given entity.
func ChangeSecret(server string, port int, clientID string, serviceID string, entity string, secret string, modEntity string, modSecret string) (string, error) {
	// e is the enity requesting the change.
	e := new(pb.Entity)
	e.ID = &entity
	e.Secret = &secret

	// me is the entity to be modified
	me := new(pb.Entity)
	me.ID = &modEntity
	me.Secret = &modSecret

	request := new(pb.ModEntityRequest)
	request.Entity = e
	request.ModEntity = me

	c, err := newClient(server, port)
	if err != nil {
		return "", err
	}

	result, err := c.ChangeSecret(context.Background(), request)
	if err != nil {
		return "", err
	}

	return result.GetMsg(), nil
}

// GroupMembers returns a list of members in a group from the server.
func GroupMembers(serverAddr string, serverPort int, clientID, serviceID, groupID string) (*pb.EntityList, error) {
	group := new(pb.Group)
	group.Name = &groupID

	request := new(pb.GroupMemberRequest)
	request.Group = group
	request.ServiceID = &serviceID
	request.ClientID = &clientID

	c, err := newClient(serverAddr, serverPort)
	if err != nil {
		return nil, err
	}

	result, err := c.ListGroupMembers(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// EntityInfo makes a request to retrieve all that is known about an
// entity that can be retrieved.  This is necessarily the entity
// object itself, and so this function returns a pointer to the entity
// protobuf.  This could be considered a leak of the type outside the
// client, but it is the most straightforward way to obtain all that
// can be obtained at once.
func EntityInfo(serverAddr string, serverPort int, clientID, serviceID, ID string) (*pb.Entity, error) {
	e := new(pb.Entity)
	e.ID = &ID

	request := new(pb.NetAuthRequest)
	request.Entity = e
	request.ServiceID = ensureServiceID(serviceID)
	request.ClientID = ensureClientID(clientID)

	c, err := newClient(serverAddr, serverPort)
	if err != nil {
		return nil, err
	}

	entity, err := c.EntityInfo(context.Background(), request)
	if err != nil {
		return nil, err
	}

	return entity, nil
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
