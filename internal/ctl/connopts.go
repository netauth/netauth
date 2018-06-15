package ctl

import (
	"flag"
	"fmt"
	"os"

	"github.com/NetAuth/NetAuth/pkg/client"
)

var (
	cfg        *client.NACLConfig
	serverAddr = flag.String("server", getServer(), "Server Address")
	serverPort = flag.Int("port", getPort(), "Server port")
	clientID   = flag.String("client", getClientID(), "Client ID to send")
	serviceID  = flag.String("service", getServiceID(), "Service ID to send")
	entity     = flag.String("entity", "", "Entity to send in the request")
	secret     = flag.String("secret", "", "Secret to send in the request")
)

// loadConfig loads the config in.  It would have been nice to do this
// in init(), but that gets called too late
func loadConfig() {
	if cfg != nil {
		return
	}
	config := os.Getenv("NACLCONFIG")
	if config == "" {
		config = "/etc/netauth.toml"
	}

	var err error
	cfg, err = client.LoadConfig(config)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Config loading error: ", err)
		return
	}
}

// Hide the entity and secret behind functions so its easier to swap
// them out later for more complex ways to get the values.
func getEntity() string { return *entity }
func getSecret() string { return *secret }

// Hide the other defaults as well
func getServer() string {
	loadConfig()
	if cfg == nil || cfg.Server == "" {
		return "localhost"
	}
	return cfg.Server
}

func getPort() int {
	loadConfig()
	if cfg == nil || cfg.Port == 0 {
		return 8080
	}
	return cfg.Port
}

func getServiceID() string {
	loadConfig()
	if cfg == nil || cfg.ServiceID == "" {
		return "netauthctl"
	}
	return cfg.ServiceID
}

func getClientID() string {
	loadConfig()
	if cfg == nil || cfg.ClientID == "" {
		hostname, err := os.Hostname()
		if err != nil {
			return ""
		}
		return hostname
	}
	return cfg.ClientID
}

// getClient attempts to return a client to the caller that is
// configured and ready to use.
func getClient() (*client.NetAuthClient, error) {
	cconf := client.NACLConfig{
		Server:    *serverAddr,
		Port:      *serverPort,
		ClientID:  *clientID,
		ServiceID: *serviceID,
	}
	return client.New(&cconf)
}
