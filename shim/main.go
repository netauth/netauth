package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"

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
	clientID = flag.String("client_id", "", "ID of the client")
	serviceID = flag.String("service_id", "authshim", "Service ID of the request")
	debug        = flag.Bool("debug", false, "Print sensitive information to aid in debugging")
)

func init() {
	// Init will setup the entity structure based on the values
	// obtained by flags, or by a stdio reader if values were not
	// set.
	flag.Parse()

	entity.ID = proto.String(*id)
	entity.Secret = proto.String(*secret)

	if entity.GetID() == "" {
		fmt.Print("Entity: ")
		reader := bufio.NewReader(os.Stdin)
		id, _ := reader.ReadString('\n')
		entity.ID = proto.String(id)
	}
	if entity.GetSecret() == "" && *authenticate {
		fmt.Print("Secret: ")
		reader := bufio.NewReader(os.Stdin)
		secret, _ := reader.ReadString('\n')
		entity.Secret = proto.String(secret)
	}
}

func main() {

	// Setup the connection and defer the close.
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", *serverAddr, *serverPort), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to NetAuth: %s", err)
	}
	// This success message is very misleading, in theory if you
	// see this then a valid connection has been made to the
	// server, but this isn't really the case.  This will show
	// connected in all cases where the Dial() function returns
	// successfully, whether or not it has actually connected to
	// the NetAuth service is another matter entirely.  In theory
	// we could fire a PingRequest() before printing this message,
	// but that's somewhat superfluous when this will all fail out
	// within the next second if there's a problem.
	log.Printf("Connected to NetAuth server at %s:%d", *serverAddr, *serverPort)
	defer conn.Close()

	// Create a client to use later on.
	client := pb.NewSystemAuthClient(conn)

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
		pingResult, err := client.Ping(context.Background(), &pb.PingRequest{})
		if err != nil {
			log.Fatalf("Ping failed: %s", err)
		}
		log.Printf("Ping successful?")
		log.Printf("%s", pingResult)
	}
}
