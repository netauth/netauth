package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/NetAuth/NetAuth/internal/server/db"
	_ "github.com/NetAuth/NetAuth/internal/server/db/impl"
	"github.com/NetAuth/NetAuth/internal/server/entity_manager"
	"github.com/NetAuth/NetAuth/internal/server/health"
	"github.com/NetAuth/NetAuth/internal/server/rpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/NetAuth/NetAuth/pkg/proto"
)

var (
	bindPort  = flag.Int("port", 8080, "Serving port, defaults to 8080")
	bindAddr  = flag.String("bind", "localhost", "Bind address, defaults to localhost")
	useTLS    = flag.Bool("tls", false, "Enable TLS, off by default")
	certFile  = flag.String("cert_file", "", "Path to certificate file")
	keyFile   = flag.String("key_file", "", "Path to key file")
	bootstrap = flag.String("make_bootstrap", "", "ID:secret to give GLOBAL_ROOT - for bootstrapping")
	db_impl   = flag.String("db", "MemDB", "Database implementation to use.")
)

func newServer() *rpc.NetAuthServer {
	// Need to setup the EMDiskInterface to pass to the entity_manager.
	db, err := db.New(*db_impl)
	if err != nil {
		log.Fatalf("Fatal database error! (%s)", err)
	}

	return &rpc.NetAuthServer{
		EM: entity_manager.New(&db),
	}
}

func main() {
	flag.Parse()

	log.Println("NetAuth server is starting!")

	// Bind early so that if this fails we can just bail out.
	sock, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *bindAddr, *bindPort))
	if err != nil {
		log.Fatalf("Could not bind! %v", err)
	}
	log.Printf("Server bound on %s:%d", *bindAddr, *bindPort)

	// Setup the TLS parameters if necessary.
	var opts []grpc.ServerOption
	if *useTLS {
		log.Printf("TLS with the certificate %s and key %s", *certFile, *keyFile)
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("TLS credentials could not be generated! %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	if !*useTLS {
		// Not using TLS in an auth server?  For shame...
		log.Println("Launching without TLS! Your passwords will be shipped in the clear!")
		log.Println("You should really start the server with -tls -key_file <keyfile> -cert_file <certfile>")
	}

	// Spit out what backends we know about
	log.Printf("The following DB backends are registered:")
	for _, b := range db.GetBackendList() {
		log.Printf("  %s", b)
	}

	// Init the new server instance
	srv := newServer()

	// Attempt to bootstrap a superuser
	if len(*bootstrap) != 0 {
		log.Println("Commencing Bootstrap")
		eParts := strings.Split(*bootstrap, ":")
		srv.EM.MakeBootstrap(eParts[0], eParts[1])
		log.Println("Bootstrap phase complete")
	}

	// If it wasn't used make sure its disabled since it can
	// create arbitrary root users.
	srv.EM.DisableBootstrap()

	// At this point the server should be ready to serve.
	health.SetGood()

	// Instantiate and launch.  This will block and the server
	// will server forever.
	log.Println("Server is launching...")
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterNetAuthServer(grpcServer, srv)
	grpcServer.Serve(sock)
}
