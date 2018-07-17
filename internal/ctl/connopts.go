package ctl

import (
	"flag"
	"fmt"
	"os"
	"os/user"

	"github.com/bgentry/speakeasy"

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
	insecure   = flag.Bool("PWN_ME", false, "Run without server verification")
	serverCert = flag.String("certificate", getServerCert(), "Certificate of the NetAuth server")
)

// loadConfig loads the config in.  It would have been nice to do this
// in init(), but that gets called too late
func loadConfig() {
	if cfg != nil {
		return
	}
	var err error
	cfg, err = client.LoadConfig("")
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Config loading error: ", err)
		return
	}
}

// Try to return the entity from the system unless overridden.
func getEntity() string {
	if *entity != "" {
		return *entity
	}
	user, err := user.Current()
	if err != nil {
		return ""
	}
	return user.Username
}

// Prompt for the secret if it wasn't provided in cleartext.
func getSecret() string {
	if *secret != "" {
		return *secret
	}
	password, err := speakeasy.Ask("Secret: ")
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	return password
}

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

func getServerCert() string {
	loadConfig()
	if cfg == nil || cfg.ServerCert == "" {
		return "/etc/netauth.cert"
	}
	return cfg.ServerCert
}

// getClient attempts to return a client to the caller that is
// configured and ready to use.
func getClient() (*client.NetAuthClient, error) {
	cconf := client.NACLConfig{
		Server:         *serverAddr,
		Port:           *serverPort,
		ClientID:       *clientID,
		ServiceID:      *serviceID,
		ServerCert:     *serverCert,
		WildlyInsecure: *insecure,
	}
	return client.New(&cconf)
}
