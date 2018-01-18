package ctl

var (
	serverAddr string
	serverPort int
	clientID   string
	serviceID  string
	entity     string
	secret     string
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
