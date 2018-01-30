package client

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/grpc"

	pb "github.com/NetAuth/NetAuth/pkg/proto"
)

type netAuthClient struct {
	c         pb.NetAuthClient
	serviceID *string
	clientID  *string
}

// New takes in the values that set up a client and builds a
// client.netAuthClient struct on which all other methods are bound.
// This drastically simplifies the construction of other functions.
func New(server string, port int, serviceID string, clientID string) (*netAuthClient, error) {
	// Setup the connection and defer the close.
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", server, port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	// Create a client to use later on.
	client := netAuthClient{
		c:         pb.NewNetAuthClient(conn),
		serviceID: ensureServiceID(serviceID),
		clientID:  ensureClientID(clientID),
	}

	return &client, nil
}

// Ping very simply pings the server.  The reply will contain the
// health status of the server as a server that replies and a server
// that can serve are two very different things (data might be
// reloading during the request).
func (n *netAuthClient) Ping() (string, error) {
	request := new(pb.PingRequest)
	request.ClientID = n.clientID

	pingResult, err := n.c.Ping(context.Background(), request)
	if err != nil {
		return "RPC Error", err
	}
	return pingResult.GetMsg(), nil
}

// Authenticate takes in an entity and a secret and tries to validate
// that the identity is legitimate by verifying the secret provided.
func (n *netAuthClient) Authenticate(entity string, secret string) (string, error) {
	e := new(pb.Entity)
	e.ID = &entity
	e.Secret = &secret

	request := new(pb.NetAuthRequest)
	request.Entity = e
	request.ClientID = n.clientID
	request.ServiceID = n.serviceID

	authResult, err := n.c.AuthEntity(context.Background(), request)
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
func (n *netAuthClient) NewEntity(entity, secret, newEntity string, newUIDNumber int32, newSecret string) (string, error) {
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

	result, err := n.c.NewEntity(context.Background(), request)
	if err != nil {
		return "", err
	}

	return result.GetMsg(), nil
}

// RemoveEntity makes a request to the server to remove the named
// entity.  This must be authorized by an entity which has the
// appropriate capabilities to fulfill the remove request
func (n *netAuthClient) RemoveEntity(entity, secret, delEntity string) (string, error) {
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

	result, err := n.c.RemoveEntity(context.Background(), request)
	if err != nil {
		return "", err
	}

	return result.GetMsg(), nil
}

// ChangeSecret changes the secret on the given modEntity under the
// authority of the given entity.
func (n *netAuthClient) ChangeSecret(entity, secret, modEntity, modSecret string) (string, error) {
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

	result, err := n.c.ChangeSecret(context.Background(), request)
	if err != nil {
		return "", err
	}

	return result.GetMsg(), nil
}

// GroupMembers returns a list of members in a group from the server.
func (n *netAuthClient) GroupMembers(groupID string) (*pb.EntityList, error) {
	group := new(pb.Group)
	group.Name = &groupID

	request := new(pb.GroupMemberRequest)
	request.Group = group
	request.ServiceID = n.serviceID
	request.ClientID = n.clientID

	result, err := n.c.ListGroupMembers(context.Background(), request)
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
func (n *netAuthClient) EntityInfo(ID string) (*pb.Entity, error) {
	e := new(pb.Entity)
	e.ID = &ID

	request := new(pb.NetAuthRequest)
	request.Entity = e
	request.ServiceID = n.serviceID
	request.ClientID = n.clientID

	entity, err := n.c.EntityInfo(context.Background(), request)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

// ModifyEntityMeta takes in a set of credentials to authorize the
// change, an entity to make the change against, and a new EntityMeta
// to apply.
func (n *netAuthClient) ModifyEntityMeta(entity, secret, modID string, meta *pb.EntityMeta) (string, error) {
	e := new(pb.Entity)
	e.ID = &entity
	e.Secret = &secret

	me := new(pb.Entity)
	me.ID = &modID
	me.Meta = meta

	request := new(pb.ModEntityRequest)
	request.Entity = e
	request.ModEntity = me

	result, err := n.c.ModifyEntityMeta(context.Background(), request)
	if err != nil {
		return "", err
	}

	return result.GetMsg(), nil
}

// NewGroup takes in a set of credentials to authorize the change and
// parameters to create a new group with.
func (n *netAuthClient) NewGroup(entity, secret, name, displayName string, gidNumber int32) (string, error) {
	e := new(pb.Entity)
	e.ID = &entity
	e.Secret = &secret

	g := new(pb.Group)
	g.Name = &name
	g.DisplayName = &displayName
	g.GidNumber = &gidNumber

	request := new(pb.ModGroupRequest)
	request.Entity = e
	request.Group = g

	result, err := n.c.NewGroup(context.Background(), request)
	if err != nil {
		return "", err
	}

	return result.GetMsg(), nil
}

// DeleteGroup deletes a group with the authorization of an existing
// entity.
func (n *netAuthClient) DeleteGroup(entity, secret, name string) (string, error) {
	e := new(pb.Entity)
	e.ID = &entity
	e.Secret = &secret

	g := new(pb.Group)
	g.Name = &name

	request := new(pb.ModGroupRequest)
	request.Entity = e
	request.Group = g

	result, err := n.c.DeleteGroup(context.Background(), request)
	if err != nil {
		return "", err
	}

	return result.GetMsg(), nil
}

// ModifyGroupMeta modifies a group with the authorization of an
// existing entity.
func (n *netAuthClient) ModifyGroupMeta(entity, secret string, update *pb.Group) (string, error) {
	e := new(pb.Entity)
	e.ID = &entity
	e.Secret = &secret

	request := new(pb.ModGroupRequest)
	request.Entity = e
	request.Group = update

	result, err := n.c.ModifyGroupMeta(context.Background(), request)
	if err != nil {
		return "", err
	}
	return result.GetMsg(), nil
}

// ListGroups lists  the groups that  are known to the  server.  These
// are just strings,  so additional requests are needed  to do things,
// but this gives some idea of what the server knows about.
func (n *netAuthClient) ListGroups() ([]*pb.Group, error) {
	request := new(pb.GroupListRequest)
	result, err := n.c.ListGroups(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return result.GetGroups(), nil
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
