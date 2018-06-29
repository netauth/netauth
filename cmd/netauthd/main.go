package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/NetAuth/NetAuth/internal/crypto"
	_ "github.com/NetAuth/NetAuth/internal/crypto/impl"
	"github.com/NetAuth/NetAuth/internal/db"
	_ "github.com/NetAuth/NetAuth/internal/db/impl"
	"github.com/NetAuth/NetAuth/internal/token"
	_ "github.com/NetAuth/NetAuth/internal/token/impl"

	"github.com/NetAuth/NetAuth/internal/health"
	"github.com/NetAuth/NetAuth/internal/rpc"
	"github.com/NetAuth/NetAuth/internal/tree"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/NetAuth/Protocol"
)

var (
	bindPort    = flag.Int("port", 8080, "Serving port, defaults to 8080")
	bindAddr    = flag.String("bind", "localhost", "Bind address, defaults to localhost")
	insecure    = flag.Bool("PWN_ME", false, "Disable TLS; Don't set on a production server!")
	certFile    = flag.String("cert_file", "netauth.cert", "Path to certificate file")
	keyFile     = flag.String("key_file", "netauth.certkey", "Path to key file")
	bootstrap   = flag.String("make_bootstrap", "", "ID:secret to give GLOBAL_ROOT - for bootstrapping")
	db_impl     = flag.String("db", "ProtoDB", "Database implementation to use.")
	crypto_impl = flag.String("crypto", "bcrypt", "Crypto implementation to use.")
)

func newServer() *rpc.NetAuthServer {
	// Need to setup the Database for use with the entity tree
	db, err := db.New(*db_impl)
	if err != nil {
		log.Fatalf("Fatal database error! (%s)", err)
	}

	crypto, err := crypto.New(*crypto_impl)
	if err != nil {
		log.Fatalf("Fatal crypto error! (%s)", err)
	}

	// Initialize the entity tree
	log.Printf("Initializing new Entity Tree with %s and %s", *db_impl, *crypto_impl)
	tree := tree.New(db, crypto)

	// Initialize the token service
	log.Println("Initializing token service")
	tokenService, err := token.New()
	if err != nil {
		log.Fatalf("Fatal error initializing token service: %s", err)
	}

	return &rpc.NetAuthServer{
		Tree:  tree,
		Token: tokenService,
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
	if !*insecure {
		log.Printf("TLS with the certificate %s and key %s", *certFile, *keyFile)
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("TLS credentials could not be loaded! %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	} else {
		// Not using TLS in an auth server?  For shame...
		log.Println("Launching without TLS! Your passwords will be shipped in the clear!")
		log.Println("You should really start the server with -tls -key_file <keyfile> -cert_file <certfile>")
	}

	// Spit out what backends we know about
	log.Printf("The following DB backends are registered:")
	for _, b := range db.GetBackendList() {
		log.Printf("  %s", b)
	}

	// Spit out what crypto backends we know about
	log.Printf("The following crypto implementations are registered:")
	for _, b := range crypto.GetBackendList() {
		log.Printf("  %s", b)
	}

	// Spit out the token services we know about
	log.Printf("The following token services are registered:")
	for _, b := range token.GetBackendList() {
		log.Printf("  %s", b)
	}

	// Init the new server instance
	srv := newServer()

	// Attempt to bootstrap a superuser
	if len(*bootstrap) != 0 {
		log.Println("Commencing Bootstrap")
		eParts := strings.Split(*bootstrap, ":")
		srv.Tree.MakeBootstrap(eParts[0], eParts[1])
		log.Println("Bootstrap phase complete")
	}

	// If it wasn't used make sure its disabled since it can
	// create arbitrary root users.
	srv.Tree.DisableBootstrap()

	// Instantiate and launch.  This will block and the server
	// will server forever.
	log.Println("Ready to Serve...")
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterNetAuthServer(grpcServer, srv)

	// Flip the status okay and launch into the RPC handling
	// phase.
	health.SetGood()
	grpcServer.Serve(sock)
}
