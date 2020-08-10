package client

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/netauth/netauth/internal/token"
	// Register the token services on import
	_ "github.com/netauth/netauth/internal/token/all"

	"github.com/netauth/netauth/internal/startup"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/netauth/protocol"
)

var (
	// ErrConfigError is returned when the configuration was
	// loaded but was missing required values.
	ErrConfigError = errors.New("required configuration values are missing")
)

func init() {
	viper.SetDefault("tls.certificate", "keys/tls.pem")
	viper.SetDefault("core.port", 1729)
	viper.SetDefault("token.backend", "jwt-rsa")
}

// New returns a complete client ready to use.
func New() (*NetAuthClient, error) {
	log := hclog.L().Named("nacl")
	token.SetParentLogger(log)

	// Logging and config are available, run deferred startup
	// hooks.
	startup.DoCallbacks()

	// Set defaults for the client ID and service ID
	hn, err := os.Hostname()
	if err != nil {
		viper.SetDefault("client.ID", "BOGUS_CLIENT")
	} else {
		viper.SetDefault("client.ID", hn)
	}
	viper.SetDefault("client.ServiceName", "BOGUS_SERVICE")

	// Setup the connection.
	conn, err := connect(false)
	if err != nil {
		log.Error("Error during connect", "error", err)
		return nil, err
	}

	// Get a tokenstore
	t, err := getTokenStore()
	if err != nil {
		// Log the error, but as there are many queries done
		// in read only mode, don't fail on it.
		log.Warn("Token storage are unavailable", "error", err)
	}

	// Get a token service, don't be a fatal error as most queries
	// don't require authentication anyway.
	ts, err := token.New(viper.GetString("token.backend"))
	if err != nil {
		log.Warn("Token validation will be unvavailable", "error", err)
	}

	// Create a client to use later on.  The value of the readonly
	// flag is set here based on if the master and server
	// addresses are the same.
	client := NetAuthClient{
		c:            pb.NewNetAuthClient(conn),
		tokenStore:   t,
		tokenService: ts,
		readonly:     viper.GetString("core.server") != viper.GetString("core.master"),
		log:          log,
	}

	return &client, nil
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

func (n *NetAuthClient) makeWritable() error {
	// If the master server is the one that we would already be
	// connected to, then just return.  Also return if we are
	// already not readonly.
	if viper.GetString("core.server") == viper.GetString("core.master") || !n.readonly {
		return nil
	}

	conn, err := connect(true)
	if err != nil {
		return err
	}
	n.c = pb.NewNetAuthClient(conn)
	return nil
}
