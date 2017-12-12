package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/NetAuth/NetAuth/pkg/client"

	"github.com/golang/protobuf/proto"

	pb "github.com/NetAuth/NetAuth/proto"
)

var (
	entity       = new(pb.Entity)
	serverAddr   = flag.String("server", "localhost", "NetAuth server to contact")
	serverPort   = flag.Int("port", 8080, "Port for the NetAuth Server")
	id           = flag.String("entity_id", "", "Entity to send")
	secret       = flag.String("entity_secret", "", "Entity secret to send")
	authenticate = flag.Bool("auth", true, "Try to authenticate the entity")
	getInfo      = flag.Bool("info", false, "Try to get info about the entity")
	ping         = flag.Bool("ping", false, "Ping the server")
	clientID     = flag.String("client_id", "", "ID of the client")
	serviceID    = flag.String("service_id", "authshim", "Service ID of the request")
	debug        = flag.Bool("debug", false, "Print sensitive information to aid in debugging")
)

func init() {
	// Init will setup the entity structure based on the values
	// obtained by flags, or by a stdio reader if values were not
	// set.
	flag.Parse()

	entity.ID = proto.String(*id)
	entity.Secret = proto.String(*secret)
}

func main() {
	// Get a client.  References to the underlying connection will
	// be cleaned up as the application tears down.
	client := client.NewClient(*serverAddr, *serverPort)

	// Grab a request object.  Its a few extra bytes that might
	// not be used, but the PingRequest case is not common so
	// premature optimization is not really worth it here.
	request := new(pb.NetAuthRequest)

	// If this is a request that will use the NetAuthRequest proto
	// then we need to populate some basic fields that those need.
	// Strictly if this fails it not critical, but it will make
	// debugging substantially harder.
	if *authenticate || *getInfo {
		if *clientID == "" {
			hostname, err := os.Hostname()
			if err != nil {
				hostname = "BOGUS_CLIENT"
			}
			request.ClientID = proto.String(hostname)
		} else {
			request.ClientID = proto.String(*clientID)
		}
		request.ServiceID = serviceID
	}

	// If we're going to authenticate we'll go ahead and setup the
	// calls and fire the request.  The case of a successful
	// authenticate request will return 0 status code, failure to
	// authenticate will return a non-zero status code.
	if *authenticate {
		log.Printf("Trying to authenticate %s", entity.GetID())
		request.Entity = entity
		authResult, err := client.AuthEntity(context.Background(), request)
		if err != nil {
			log.Fatalf("Could not auth: %s", err)
		}
		if authResult == nil {
			log.Fatal("recieved nil reply for AuthEntity()")
		}
		log.Printf("%v", authResult)
	}

	// Its desireable to ping the server during testing to find
	// things that might go wrong.  This function also provides a
	// means to check the server's health by monitoring systems.
	if *ping {
		log.Printf("Pinging the server")
		request := new(pb.PingRequest)
		if *clientID == "" {
			hostname, err := os.Hostname()
			if err != nil {
				hostname = "BOGUS_CLIENT"
			}
			request.ClientID = proto.String(hostname)
		} else {
			request.ClientID = proto.String(*clientID)
		}
		pingResult, err := client.Ping(context.Background(), request)
		if err != nil {
			log.Fatalf("Ping failed: %s", err)
		}
		log.Printf("%s", pingResult)
	}
}
