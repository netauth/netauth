package netauth

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/netauth/netauth/internal/token"

	// The default token service is the jwt implementation, and
	// since its internal, the client needs to import it on behalf
	// of consumers.
	_ "github.com/netauth/netauth/internal/token/jwt"

	rpc "github.com/netauth/protocol/v2"
)

func init() {
	viper.SetDefault("tls.certificate", "keys/tls.pem")
	viper.SetDefault("core.port", 1729)
}

// New returns a client initialized, connected, and ready to use.
func New() (*Client, error) {
	l := hclog.L().Named("cli")

	conn, err := connect(false)
	if err != nil {
		return nil, err
	}

	cache, err := NewTokenCache(viper.GetString("token.cache"))
	if err != nil {
		return nil, err
	}

	ts, err := token.New()
	if err != nil {
		l.Warn("Token service initialization error", "error", err)
	}

	hn, err := os.Hostname()
	if err != nil {
		viper.SetDefault("client.ID", "BOGUS_CLIENT")
	} else {
		viper.SetDefault("client.ID", hn)
	}

	return &Client{
		TokenCache: cache,
		Service:    ts,
		rpc:        rpc.NewNetAuth2Client(conn),
		log:        l,
		clientName: viper.GetString("client.ID"),
	}, nil
}

func connect(writable bool) (*grpc.ClientConn, error) {
	addr := viper.GetString("core.server")

	// This has to happen here since it needs to happen after
	// everything else is already parsed.
	if viper.GetString("core.master") == "" {
		viper.Set("core.master", viper.GetString("core.server"))
	}

	if writable {
		addr = viper.GetString("core.master")
	}

	var opts []grpc.DialOption
	if viper.GetBool("tls.pwn_me") {
		opts = []grpc.DialOption{grpc.WithInsecure()}
	} else {
		// If this is a relative path its relative to the home
		// directory.
		certPath := viper.GetString("tls.certificate")
		if !filepath.IsAbs(certPath) {
			certPath = filepath.Join(viper.GetString("core.home"), certPath)
		}

		creds, err := credentials.NewClientTLSFromFile(certPath, "")
		if err != nil {
			return nil, err
		}
		opts = []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	}
	return grpc.Dial(
		fmt.Sprintf("%s:%d", addr, viper.GetInt("core.port")),
		opts...,
	)
}

// SetServiceName sets the self identified service this client serves.
// This should be set prior to making any calls to the server.
func (c *Client) SetServiceName(s string) {
	c.serviceName = s
}

func (c *Client) makeWritable() error {
	// If the master server is the one that we would already be
	// connected to, then just return.  Also return if we are
	// already not readonly.
	if viper.GetString("core.server") == viper.GetString("core.master") || c.writeable {
		return nil
	}

	conn, err := connect(true)
	if err != nil {
		return err
	}
	c.rpc = rpc.NewNetAuth2Client(conn)
	c.writeable = true
	return nil
}
