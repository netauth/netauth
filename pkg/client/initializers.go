package client

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/NetAuth/NetAuth/internal/token"
	// Register the token services on import
	_ "github.com/NetAuth/NetAuth/internal/token/all"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/NetAuth/Protocol"
)

var (
	// ErrConfigError is returned when the configuration was
	// loaded but was missing required values.
	ErrConfigError = errors.New("Required configuration values are missing")
)

// New returns a complete client ready to use.
func New() (*NetAuthClient, error) {
	// Set defaults for the client ID and service ID
	hn, err := os.Hostname()
	if err != nil {
		viper.SetDefault("client.ID", "BOGUS_CLIENT")
	} else {
		viper.SetDefault("client.ID", hn)
	}
	viper.SetDefault("client.ServiceName", "BOGUS_SERVICE")

	// Setup the connection.
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
			log.Printf("Could not load certificate: %s", err)
			return nil, err
		}
		opts = []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	}
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", viper.GetString("core.server"), viper.GetInt("core.port")),
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
		tokenStore:   t,
		tokenService: ts,
	}

	return &client, nil
}
