package ctl

import (
	"fmt"
	"os"
	"os/user"

	"github.com/bgentry/speakeasy"
	"github.com/spf13/pflag"

	"github.com/NetAuth/NetAuth/pkg/client"
)

var (
	serverAddr = flag.String("server", getServer(), "Server Address")
	serverPort = flag.Int("port", getPort(), "Server port")
	clientID   = flag.String("client", getClientID(), "Client ID to send")
	serviceID  = flag.String("service", getServiceID(), "Service ID to send")
	entity     = flag.String("entity", "", "Entity to send in the request")
	secret     = flag.String("secret", "", "Secret to send in the request")
	insecure   = flag.Bool("PWN_ME", false, "Run without server verification")
	serverCert = flag.String("certificate", getServerCert(), "Certificate of the NetAuth server")
)

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
	var err error
	*secret, err = speakeasy.Ask("Secret: ")
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	return *secret
}

func getToken(c *client.NetAuthClient, entity string) (string, error) {
	t, err := c.GetToken(entity, "")
	switch err {
	case nil:
		return t, nil
	case client.ErrTokenUnavailable:
		return c.GetToken(entity, getSecret())
	default:
		return "", err
	}
}
