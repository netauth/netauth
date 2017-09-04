package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	port     = flag.Int("port", 8080, "Serving port, defaults to 8080")
	bindAddr = flag.String("bind", "localhost", "Bind address, defaults to localhost")
	useTLS   = flag.Bool("tls", false, "Enable TLS, off by default")
	certFile = flag.String("cert_file", "", "Path to certificate file")
	keyFile  = flag.String("key_file", "", "Path to key file")
)

func main() {
	flag.Parse()

	log.Println("NetAuth server is starting!")
	
	// Bind early so that if this fails we can just bail out.
	sock, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *bindAddr, *port))
	if err != nil {
		log.Fatalf("could not bind! %v", err)
	}
	log.Printf("server bound on %s:%d", *bindAddr, *port)
	
	// Setup the TLS parameters if necessary.
	var opts []grpc.ServerOption
	if *useTLS {
		log.Printf("this server will use TLS with the certificate %s and key %s", *certFile, *keyFile)
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("TLS credentials could not be generated! %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	if ! *useTLS {
		// Not using TLS in an auth server?  For shame...
		log.Println("launching without TLS! Your passwords will be shipped in the clear!")
		log.Println("You should really start the server with -tls -key_file <keyfile> -cert_file <certfile>")
	}

	// Instantiate and launch.  This will block and the server
	// will server forever.
	log.Println("Server is launching...")
	grpcServer := grpc.NewServer(opts...)
	grpcServer.Serve(sock)
}
