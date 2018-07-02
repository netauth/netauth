package client

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/NetAuth/NetAuth/internal/token"
	// Register the token services on import
	_ "github.com/NetAuth/NetAuth/internal/token/all"

	"github.com/BurntSushi/toml"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/NetAuth/Protocol"
)

var (
	// ErrConfigError is returned when the configuration was
	// loaded but was missing required values.
	ErrConfigError = errors.New("Required configuration values are missing")
)

// New takes in a NACLConfig pointer and uses this to bootstrap a
// client.  If the pointer is nil, then the config will be loaded from
// disk from the default location.
func New(cfg *NACLConfig) (*NetAuthClient, error) {
	if cfg == nil {
		// Load from disk
		var err error
		cfg, err = LoadConfig("")
		if err != nil {
			return nil, err
		}
	}
	cfg.ServiceID = ensureServiceID(cfg.ServiceID)
	cfg.ClientID = ensureClientID(cfg.ClientID)

	// Make sure the server/port tuple is defined.
	if cfg.Server == "" {
		return nil, ErrConfigError
	}
	if cfg.Port == 0 {
		cfg.Port = 8080
	}

	// Setup the connection.
	var opts []grpc.DialOption
	if cfg.WildlyInsecure {
		opts = []grpc.DialOption{grpc.WithInsecure()}
	} else {
		// If it wasn't set then pull it from the default
		// location and hope its there.
		if cfg.ServerCert == "" {
			cfg.ServerCert = "/etc/netauth.cert"
		}

		creds, err := credentials.NewClientTLSFromFile(cfg.ServerCert, "")
		if err != nil {
			log.Printf("Could not load certificate: %s", err)
			return nil, err
		}
		opts = []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	}
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.Server, cfg.Port),
		opts...,
	)
	if err != nil {
		return nil, err
	}

	// Get a tokenstore
	t, err := getTokenStore()
	if err != nil {
		// Log the error, but as there are many queries done
		// in read only mode, don't fail on it.
		log.Println(err)
	}

	// Get a token service, don't be a fatal error as most queries
	// don't require authentication anyway.
	ts, err := token.New()
	if err != nil {
		log.Println(err)
	}

	// Create a client to use later on.
	client := NetAuthClient{
		c:            pb.NewNetAuthClient(conn),
		cfg:          cfg,
		tokenStore:   t,
		tokenService: ts,
	}

	return &client, nil
}

// LoadConfig fetches the configuration file from disk in the default
// location, or the provided path if specified.
func LoadConfig(cfgpath string) (*NACLConfig, error) {
	if cfgpath == "" {
		cfgpath = os.Getenv("NACLCONFIG")
		if cfgpath == "" {
			// If it wasn't set, this is the location to
			// load from.  At some point this path should
			// come about via an OS agnostic way since
			// Windows doesn't have an /etc to load from.
			cfgpath = "/etc/netauth.toml"
		}
	}

	// Actually load the config
	var cfg NACLConfig
	_, err := toml.DecodeFile(cfgpath, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
