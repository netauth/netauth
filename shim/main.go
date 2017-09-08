package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	pb "github.com/NetAuth/NetAuth/proto"

	"google.golang.org/grpc"
)

var (
	serverAddr   = flag.String("server", "localhost", "NetAuth server to contact (localhost)")
	serverPort   = flag.Int("port", 8080, "Port for the NetAuth Server (8080)")
	entityID     = flag.String("entity_id", "", "Entity to send")
	entitySecret = flag.String("entity_secret", "", "Entity secret to send")
	authenticate = flag.Bool("auth", true, "Try to authenticate the entity")
	getInfo      = flag.Bool("info", false, "Try to get info about the entity")
	debug        = flag.Bool("debug", false, "Print sensitive information to aid in debugging")
)

func authEntity() error {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", *serverAddr, *serverPort), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()



	return nil
}

func main() {
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)
	if *entityID == "" {
		fmt.Print("Entity: ")
		*entityID, _ = reader.ReadString('\n')
	}
	if *entitySecret == "" && *authenticate {
		fmt.Print("Secret: ")
		*entitySecret, _ = reader.ReadString('\n')
	}

	log.Printf("Entity: %s", *entityID)
	if *debug {
		log.Printf("Entity Secret: %s", *entitySecret)
	}

	conn, err := grpc.Dial(*serverAddr)
	if err != nil {
		log.Fatalf("Could not connect to NetAuth: %s", err)
	}
	defer conn.Close()

	// Create a client to use later on.
	client := pb.NewSystemAuthClient(conn)

	if *authenticate {
		log.Printf("Trying to authenticate %s to server %s:%d", *entityID, *serverAddr, *serverPort)
		err := authEntity()
		if err != nil {
			log.Fatalf("could not auth: %s", err)
		}
	}
}
