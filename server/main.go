package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/golang/protobuf/proto"

	pb "github.com/NetAuth/NetAuth/proto"
)

var (
	bindPort      = flag.Int("port", 8080, "Serving port, defaults to 8080")
	bindAddr      = flag.String("bind", "localhost", "Bind address, defaults to localhost")
	useTLS        = flag.Bool("tls", false, "Enable TLS, off by default")
	certFile      = flag.String("cert_file", "", "Path to certificate file")
	keyFile       = flag.String("key_file", "", "Path to key file")
	serverHealthy = true
)

type netAuthServer struct{}

func (s *netAuthServer) AuthEntity(ctx context.Context, netAuthRequest *pb.NetAuthRequest) (*pb.AuthResult, error) {
	// This must always be defaulted to false here.  Arguably the
	// security of the entire system stems from here where this
	// starts out as false and will require a positive action
	// below to set it true.
	var success = false

	// Go ahead and say who is making this request, and from
	// where, and for what.  This is for diagnostics, and is not
	// really intended to be used for security purposes, but can
	// be nice to look at if things fail below.
	log.Printf("Authenticating %s for service %s to client %s",
		netAuthRequest.GetEntity().GetID(),
		netAuthRequest.GetServiceID(),
		netAuthRequest.GetClientID())

	// Construct and return the response.
	result := new(pb.AuthResult)
	result.Success = &success
	return result, nil
}

func (s *netAuthServer) EntityInfo(ctx context.Context, netAuthRequest *pb.NetAuthRequest) (*pb.EntityMeta, error) {
	return &pb.EntityMeta{}, nil
}

func (s *netAuthServer) Ping(ctx context.Context, pingRequest *pb.PingRequest) (*pb.PingResponse, error) {
	// Ping takes in a request from the client, and then replies
	// with a Pong containing the server status.

	log.Printf("Ping from %s", pingRequest.GetClientID())

	reply := new(pb.PingResponse)
	reply.Healthy = &serverHealthy
	hostname, err := os.Hostname()
	if err != nil {
		log.Println("Hostname could not be determined!")
		hostname = "BOGUS_HOST"
	}
	reply.Msg = proto.String(fmt.Sprintf("NetAuth server on %s is ready to serve!", hostname))
	return reply, nil
}

func newServer() *netAuthServer {
	return new(netAuthServer)
}

func main() {
	flag.Parse()

	log.Println("NetAuth server is starting!")

	// Bind early so that if this fails we can just bail out.
	sock, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *bindAddr, *bindPort))
	if err != nil {
		log.Fatalf("could not bind! %v", err)
	}
	log.Printf("server bound on %s:%d", *bindAddr, *bindPort)

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

	if !*useTLS {
		// Not using TLS in an auth server?  For shame...
		log.Println("launching without TLS! Your passwords will be shipped in the clear!")
		log.Println("You should really start the server with -tls -key_file <keyfile> -cert_file <certfile>")
	}

	// Instantiate and launch.  This will block and the server
	// will server forever.
	log.Println("Server is launching...")
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterSystemAuthServer(grpcServer, newServer())
	grpcServer.Serve(sock)
}
