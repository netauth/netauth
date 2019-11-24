package netauth

import (
	"context"

	pb "github.com/netauth/protocol"
	rpc "github.com/netauth/protocol/v2"
)

// AuthEntity performs authentication for an entity.  It does not
// perform token acquisition, so if your request will require a token,
// ensure that you have obtained one already.
func (c *Client) AuthEntity(ctx context.Context, entity, secret string) error {
	ctx = c.appendMetadata(ctx)
	r := rpc.AuthRequest{
		Entity: &pb.Entity{
			ID: &entity,
		},
		Secret: &secret,
	}
	_, err := c.rpc.AuthEntity(ctx, &r)
	return err
}

// AuthGetToken performs authentication for an entity and if
// successful will return a token which can be used to authenticate
// future requests.
func (c *Client) AuthGetToken(ctx context.Context, entity, secret string) (string, error) {
	ctx = c.appendMetadata(ctx)
	r := rpc.AuthRequest{
		Entity: &pb.Entity{
			ID: &entity,
		},
		Secret: &secret,
	}
	res, err := c.rpc.AuthGetToken(ctx, &r)
	return res.GetToken(), err
}

// AuthValidateToken performs server-side token validation.  This can
// be useful when symmetric token algorithms are in use and clients
// are unable to validate tokens locally, or if you simply don't trust
// the local validation option.
func (c *Client) AuthValidateToken(ctx context.Context, token string) error {
	ctx = c.appendMetadata(ctx)
	r := rpc.AuthRequest{
		Token: &token,
	}
	_, err := c.rpc.AuthValidateToken(ctx, &r)
	return err
}

// AuthChangeSecret changes the secret for a given entity.  If the
// entity is changing its own secret, then the original secret must be
// supplied.  If an administrator is changing the secret, an
// appropriate token must be present.
func (c *Client) AuthChangeSecret(ctx context.Context, entity, secret, oldsecret string) error {
	if err := c.makeWritable(); err != nil {
		return err
	}

	ctx = c.appendMetadata(ctx)
	r := rpc.AuthRequest{
		Entity: &pb.Entity{
			ID:     &entity,
			Secret: &oldsecret,
		},
		Secret: &secret,
	}
	_, err := c.rpc.AuthChangeSecret(ctx, &r)
	return err
}
