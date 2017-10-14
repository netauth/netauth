package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	pb "github.com/NetAuth/NetAuth/proto"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
)

var (
	entity       = new(pb.Entity)
	serverAddr   = flag.String("server", "localhost", "NetAuth server to contact")
	serverPort   = flag.Int("port", 8080, "Port for the NetAuth Server")
	id           = flag.String("entity_id", "", "Entity to send")
	secret       = flag.String("entity_secret", "", "Entity secret to send")
	authenticate = flag.Bool("auth", true, "Try to authenticate the entity")
	getInfo      = flag.Bool("info", false, "Try to get info about the entity")
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
	log.Printf("Entity: %s", entity.GetID())
	if *debug {
		log.Printf("Entity Secret: %s", entity.GetSecret())
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", *serverAddr, *serverPort), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to NetAuth: %s", err)
	}
	log.Printf("Successfully connected to NetAuth server at %s:%d", *serverAddr, *serverPort)
	defer conn.Close()

	// Create a client to use later on.
	client := pb.NewSystemAuthClient(conn)

	if *authenticate {
		log.Printf("Trying to authenticate %s", entity.GetID())
		authResult, err := client.AuthEntity(context.Background(), entity)
		if err != nil {
			log.Fatalf("Could not auth: %s", err)
		}
		if authResult == nil {
			log.Fatal("recieved nil reply for AuthEntity()")
		}
		log.Printf("%v", authResult)
	}
}
