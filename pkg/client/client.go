package client

import (
	"context"
	"strings"

	"github.com/NetAuth/NetAuth/internal/token"
	"github.com/NetAuth/NetAuth/internal/tree"

	"github.com/golang/protobuf/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/NetAuth/Protocol"
)

// The NetAuthClient is the logical abstraction on top of the gRPC
// client form the Protobuf.  This includes the additional components
// such as the TokenService and the TokenStore, as well as the config
// structures that drive the client.
type NetAuthClient struct {
	c            pb.NetAuthClient
	tokenStore   TokenStore
	tokenService token.Service
}

// Ping very simply pings the server.  The reply will contain the
// health status of the server as a server that replies and a server
// that can serve are two very different things (data might be
// reloading during the request).
func (n *NetAuthClient) Ping() (*pb.PingResponse, error) {
	request := &pb.PingRequest{
		Info: clientInfo(),
	}

	result, err := n.c.Ping(context.Background(), request)
	if status.Code(err) != codes.OK {
		return nil, err
	}
	return result, nil
}

// Authenticate takes in an entity and a secret and tries to validate
// that the identity is legitimate by verifying the secret provided.
func (n *NetAuthClient) Authenticate(entity string, secret string) (*pb.SimpleResult, error) {
	request := pb.NetAuthRequest{
		Entity: &pb.Entity{
			ID:     &entity,
			Secret: &secret,
		},
		Info: clientInfo(),
	}

	result, err := n.c.AuthEntity(context.Background(), &request)
	if status.Code(err) != codes.OK {
		return nil, err
	}
	return result, nil
}

// ValidateToken sends the token to the server for validation.  This
// is effectively asking the server to authenticate the token and not
// do anything else.  Returns a comment from the server and an error.
func (n *NetAuthClient) ValidateToken(entity string) (*pb.SimpleResult, error) {
	t, err := n.getTokenFromStore(entity)
	if err != nil {
		return nil, err
	}

	request := pb.NetAuthRequest{
		Entity: &pb.Entity{
			ID: &entity,
		},
		AuthToken: &t,
		Info:      clientInfo(),
	}

	result, err := n.c.ValidateToken(context.Background(), &request)
	if status.Code(err) != codes.OK {
		return nil, err
	}
	return result, nil
}

// ChangeSecret crafts a modEntity request with the correct fields to
// change an entity secret either via self authentication or via token
// authentication which is held by an appropriate administrator.
func (n *NetAuthClient) ChangeSecret(e, s, me, ms, t string) (*pb.SimpleResult, error) {
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
		Info:      clientInfo(),
	}

	result, err := n.c.ChangeSecret(context.Background(), &request)
	if status.Code(err) != codes.OK {
		return nil, err
	}
	return result, nil
}

// NewEntity crafts a modEntity request with the correct fields to
// create a new entity.
func (n *NetAuthClient) NewEntity(id string, uidn int32, secret, t string) (*pb.SimpleResult, error) {
	request := pb.ModEntityRequest{
		Entity: &pb.Entity{
			ID:     &id,
			Number: &uidn,
			Secret: &secret,
		},
		AuthToken: &t,
		Info:      clientInfo(),
	}

	result, err := n.c.NewEntity(context.Background(), &request)
	if status.Code(err) != codes.OK {
		return nil, err
	}
	return result, nil
}

// RemoveEntity removes an entity by the given name.  Only the
// 'entity' field of the modEntityRequest is required.
func (n *NetAuthClient) RemoveEntity(id, token string) (*pb.SimpleResult, error) {
	request := pb.ModEntityRequest{
		Entity: &pb.Entity{
			ID: &id,
		},
		AuthToken: &token,
		Info:      clientInfo(),
	}

	result, err := n.c.RemoveEntity(context.Background(), &request)
	if status.Code(err) != codes.OK {
		return nil, err
	}
	return result, nil
}

// EntityInfo btains the entity object with the secure fields
// redacted.  This is primarily used for displaying the values of the
// metadata struct internally.
func (n *NetAuthClient) EntityInfo(id string) (*pb.Entity, error) {
	request := pb.NetAuthRequest{
		Entity: &pb.Entity{
			ID: &id,
		},
		Info: clientInfo(),
	}
	result, err := n.c.EntityInfo(context.Background(), &request)
	if status.Code(err) != codes.OK {
		return nil, err
	}
	return result, nil
}

// ModifyEntityMeta makes an authenticated request to the server to
// update the metadata of an entity.
func (n *NetAuthClient) ModifyEntityMeta(id, t string, meta *pb.EntityMeta) (*pb.SimpleResult, error) {
	request := pb.ModEntityRequest{
		Entity: &pb.Entity{
			ID:   &id,
			Meta: meta,
		},
		AuthToken: &t,
		Info:      clientInfo(),
	}

	result, err := n.c.ModifyEntityMeta(context.Background(), &request)
	if status.Code(err) != codes.OK {
		return nil, err
	}
	return result, nil
}

// ModifyEntityKeys modifies the keys on an entity, this action must
// be authorized.
func (n *NetAuthClient) ModifyEntityKeys(t, e, m, kt, kv string) ([]string, error) {
	request := pb.ModEntityKeyRequest{
		Entity: &pb.Entity{
			ID: &e,
		},
		AuthToken: &t,
		Mode:      &m,
		Type:      &kt,
		Key:       &kv,
		Info:      clientInfo(),
	}
	result, err := n.c.ModifyEntityKeys(context.Background(), &request)
	if status.Code(err) != codes.OK {
		return nil, err
	}

	keys := []string{}
	for _, k := range result.GetKeys() {
		parts := strings.Split(k, ":")
		if parts[0] == kt {
			keys = append(keys, parts[1])
		}
	}
	return keys, nil
}

// ModifyUntypedEntityMeta manages actions on the untyped metadata
// storage.
func (n *NetAuthClient) ModifyUntypedEntityMeta(t, e, m, k, v string) (map[string]string, error) {
	request := pb.ModEntityMetaRequest{
		Entity: &pb.Entity{
			ID: &e,
		},
		AuthToken: &t,
		Mode:      &m,
		Key:       &k,
		Value:     &v,
		Info:      clientInfo(),
	}

	result, err := n.c.ModifyUntypedEntityMeta(context.Background(), &request)
	if status.Code(err) != codes.OK {
		return nil, err
	}

	utm := make(map[string]string)
	for _, kv := range result.GetUntypedMeta() {
		parts := strings.SplitN(kv, ":", 2)
		utm[parts[0]] = parts[1]
	}
	return utm, nil
}

// NewGroup creates a new group with the given name, display name, and
// group number.  This action must be authorized.
func (n *NetAuthClient) NewGroup(name, displayname, managedby, t string, number int) (*pb.SimpleResult, error) {
	gid := int32(number)
	request := pb.ModGroupRequest{
		Group: &pb.Group{
			Name:        &name,
			DisplayName: &displayname,
			Number:      &gid,
			ManagedBy:   &managedby,
		},
		AuthToken: &t,
		Info:      clientInfo(),
	}

	result, err := n.c.NewGroup(context.Background(), &request)
	if status.Code(err) != codes.OK {
		return nil, err
	}
	return result, nil
}

// DeleteGroup removes a group by name.  This action must be
// authorized.
func (n *NetAuthClient) DeleteGroup(name, t string) (*pb.SimpleResult, error) {
	request := pb.ModGroupRequest{
		Group: &pb.Group{
			Name: &name,
		},
		AuthToken: &t,
		Info:      clientInfo(),
	}

	result, err := n.c.DeleteGroup(context.Background(), &request)
	if status.Code(err) != codes.OK {
		return nil, err
	}
	return result, nil
}

// ListGroups returns a list of groups to the caller.  This action
// does not require authorization.
func (n *NetAuthClient) ListGroups(entity string, indirects bool) ([]*pb.Group, error) {
	request := pb.GroupListRequest{
		Info:             clientInfo(),
		IncludeIndirects: &indirects,
	}

	if entity != "" {
		request.Entity = &pb.Entity{ID: &entity}
	}

	result, err := n.c.ListGroups(context.Background(), &request)
	if status.Code(err) != codes.OK {
		return nil, err
	}
	return result.GetGroups(), nil
}

// GroupInfo provides information about a single group.
func (n *NetAuthClient) GroupInfo(name string) (*pb.GroupInfoResult, error) {
	request := pb.ModGroupRequest{
		Info: clientInfo(),
		Group: &pb.Group{
			Name: &name,
		},
		AuthToken: proto.String(""),
	}

	result, err := n.c.GroupInfo(context.Background(), &request)
	if status.Code(err) != codes.OK {
		return nil, err
	}
	return result, nil
}

// ModifyGroupMeta allows a group's metadata to be altered after the
// fact.  This action must be authorized.
func (n *NetAuthClient) ModifyGroupMeta(group *pb.Group, token string) (*pb.SimpleResult, error) {
	request := pb.ModGroupRequest{
		Group:     group,
		AuthToken: &token,
		Info:      clientInfo(),
	}

	result, err := n.c.ModifyGroupMeta(context.Background(), &request)
	if status.Code(err) != codes.OK {
		return nil, err
	}
	return result, nil
}

// AddEntityToGroup modifies direct membership of entities.  This
// action must be authorized.
func (n *NetAuthClient) AddEntityToGroup(t, g, e string) (*pb.SimpleResult, error) {
	request := pb.ModEntityMembershipRequest{
		Entity: &pb.Entity{
			ID: &e,
		},
		Group: &pb.Group{
			Name: &g,
		},
		AuthToken: &t,
		Info:      clientInfo(),
	}

	result, err := n.c.AddEntityToGroup(context.Background(), &request)
	if status.Code(err) != codes.OK {
		return nil, err
	}
	return result, nil
}

// RemoveEntityFromGroup modifies direct membership of entities.  This
// action must be authorized.
func (n *NetAuthClient) RemoveEntityFromGroup(t, g, e string) (*pb.SimpleResult, error) {
	request := pb.ModEntityMembershipRequest{
		Entity: &pb.Entity{
			ID: &e,
		},
		Group: &pb.Group{
			Name: &g,
		},
		AuthToken: &t,
		Info:      clientInfo(),
	}

	result, err := n.c.RemoveEntityFromGroup(context.Background(), &request)
	if status.Code(err) != codes.OK {
		return nil, err
	}
	return result, nil
}

// ListGroupMembers returns a list of members for the requested group.
// This action does not require authorization.
func (n *NetAuthClient) ListGroupMembers(g string) ([]*pb.Entity, error) {
	request := pb.GroupMemberRequest{
		Group: &pb.Group{
			Name: &g,
		},
		Info: clientInfo(),
	}

	result, err := n.c.ListGroupMembers(context.Background(), &request)
	if status.Code(err) != codes.OK {
		return nil, err
	}
	return result.GetMembers(), nil
}

// ModifyGroupExpansions modifies the parent/child status of the provided groups.
// This action must be authorized.
func (n *NetAuthClient) ModifyGroupExpansions(t, p, c, m string) (*pb.SimpleResult, error) {
	m = strings.ToUpper(m)
	mode := pb.ExpansionMode(pb.ExpansionMode_value[m])

	request := pb.ModGroupNestingRequest{
		Info:      clientInfo(),
		AuthToken: &t,
		ParentGroup: &pb.Group{
			Name: &p,
		},
		ChildGroup: &pb.Group{
			Name: &c,
		},
		Mode: &mode,
	}

	result, err := n.c.ModifyGroupNesting(context.Background(), &request)
	if status.Code(err) != codes.OK {
		return nil, err
	}
	return result, nil
}

// ModifyUntypedGroupMeta manages actions on the untyped metadata
// storage.
func (n *NetAuthClient) ModifyUntypedGroupMeta(t, g, m, k, v string) (map[string]string, error) {
	request := pb.ModGroupMetaRequest{
		Group: &pb.Group{
			Name: &g,
		},
		AuthToken: &t,
		Mode:      &m,
		Key:       &k,
		Value:     &v,
		Info:      clientInfo(),
	}

	result, err := n.c.ModifyUntypedGroupMeta(context.Background(), &request)
	if status.Code(err) != codes.OK {
		return nil, err
	}

	utm := make(map[string]string)
	for _, kv := range result.GetUntypedMeta() {
		parts := strings.SplitN(kv, ":", 2)
		utm[parts[0]] = parts[1]
	}
	return utm, nil
}

// ManageCapabilities modifies the capabilities present on an entity
// or group.  This action must be authorized.
func (n *NetAuthClient) ManageCapabilities(t, e, g, c, m string) (*pb.SimpleResult, error) {
	capID, ok := pb.Capability_value[c]
	if !ok {
		return nil, tree.ErrUnknownCapability
	}
	cap := pb.Capability(capID)

	request := pb.ModCapabilityRequest{
		Info:       clientInfo(),
		AuthToken:  &t,
		Mode:       &m,
		Capability: &cap,
	}

	if e != "" {
		request.Entity = &pb.Entity{ID: &e}
	} else if g != "" {
		request.Group = &pb.Group{Name: &g}
	}

	result, err := n.c.ManageCapabilities(context.Background(), &request)
	if status.Code(err) != codes.OK {
		return nil, err
	}
	return result, nil
}

// LockEntity locks an entity which prevents validation of an entity
// secret.
func (n *NetAuthClient) LockEntity(t, e string) (*pb.SimpleResult, error) {
	request := pb.NetAuthRequest{
		Entity: &pb.Entity{
			ID: &e,
		},
		AuthToken: &t,
		Info:      clientInfo(),
	}

	result, err := n.c.LockEntity(context.Background(), &request)
	if status.Code(err) != codes.OK {
		return nil, err
	}
	return result, nil
}

// UnlockEntity unlocks an entity which was previously locked.
func (n *NetAuthClient) UnlockEntity(t, e string) (*pb.SimpleResult, error) {
	request := pb.NetAuthRequest{
		Entity: &pb.Entity{
			ID: &e,
		},
		AuthToken: &t,
		Info:      clientInfo(),
	}

	result, err := n.c.UnlockEntity(context.Background(), &request)
	if status.Code(err) != codes.OK {
		return nil, err
	}
	return result, nil
}

// SearchEntities takes a string search expression to search for
// entites on the server.
func (n *NetAuthClient) SearchEntities(expr string) (*pb.EntityList, error) {
	request := pb.SearchRequest{
		Expression: &expr,
		Info: clientInfo(),
	}
	
	result, err := n.c.SearchEntities(context.Background(), &request)
	if status.Code(err) != codes.OK {
		return nil, err
	}
	return result, nil
}

// SearchGroups takes a string search expression to search for
// entites on the server.
func (n *NetAuthClient) SearchGroups(expr string) (*pb.GroupList, error) {
	request := pb.SearchRequest{
		Expression: &expr,
		Info: clientInfo(),
	}
	
	result, err := n.c.SearchGroups(context.Background(), &request)
	if status.Code(err) != codes.OK {
		return nil, err
	}
	return result, nil
}

func clientInfo() *pb.ClientInfo {
	i := viper.GetString("client.ID")
	s := viper.GetString("client.ServiceName")

	return &pb.ClientInfo{
		ID:      &i,
		Service: &s,
	}
}
