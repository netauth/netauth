package ctl

import (
	"os"

	"github.com/NetAuth/NetAuth/pkg/client"
	"github.com/imdario/mergo"
)

var (
	serverAddr string
	serverPort int
	clientID   string
	serviceID  string
	entity     string
	secret     string
	configpath string
)

// SetServerAddr sets the server address varaiable for the rpc options
func SetServerAddr(s string) { serverAddr = s }

// SetServerPort sets the server port variable for the rpc options
func SetServerPort(p int) { serverPort = p }

// SetClientID sets the client ID for the rpc options.  This is set by
// the client and is not to be used for security purposes.
func SetClientID(s string) { clientID = s }

// SetServiceID sets the service ID for the rpc options.  This is set
// by the client and is not to be used for security purposes.
func SetServiceID(s string) { serviceID = s }

// SetEntity sets the entity for all subcommands.
func SetEntity(s string) { entity = s }

// SetSecret sets the secret for all subcommands.
func SetSecret(s string) { secret = s }

// SetConfigPath sets the path to the library configuration file
func SetConfigPath(s string) { configpath = s }

// getClient attempts to return a client to the caller that is
// configured and ready to use.
func getClient() (*client.NetAuthClient, error) {
	override := client.NACLConfig{
		Server:    serverAddr,
		Port:      serverPort,
		ClientID:  clientID,
		ServiceID: serviceID,
	}

	cfg, err := client.LoadConfig(configpath)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if err := mergo.Merge(cfg, override, mergo.WithOverride); err != nil {
		return nil, err
	}

	return client.New(cfg)
}
