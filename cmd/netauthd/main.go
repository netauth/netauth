package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/NetAuth/NetAuth/internal/crypto"
	_ "github.com/NetAuth/NetAuth/internal/crypto/all"
	"github.com/NetAuth/NetAuth/internal/db"
	_ "github.com/NetAuth/NetAuth/internal/db/all"
	"github.com/NetAuth/NetAuth/internal/token"
	_ "github.com/NetAuth/NetAuth/internal/token/all"

	"github.com/NetAuth/NetAuth/internal/rpc"
	"github.com/NetAuth/NetAuth/internal/tree"
	_ "github.com/NetAuth/NetAuth/internal/tree/hooks"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/NetAuth/Protocol"
)

var (
	bootstrap = pflag.String("server.bootstrap", "", "ID:secret to give GLOBAL_ROOT - for bootstrapping")
	insecure  = pflag.Bool("tls.PWN_ME", false, "Disable TLS; Don't set on a production server!")

	writeDefConfig = pflag.String("write-config", "", "Write the default configuration to the specified file")
)

func init() {
	pflag.String("tls.certificate", "keys/tls.crt", "Path to certificate file")
	pflag.String("tls.key", "keys/tls.key", "Path to key file")

	pflag.String("server.bind", "localhost", "Bind address, defaults to localhost")
	pflag.Int("server.port", 1729, "Serving port")
	pflag.String("core.home", "", "Base directory for NetAuth")

	pflag.String("db.backend", "ProtoDB", "Database implementation to use")

	pflag.String("crypto.backend", "bcrypt", "Cryptography system to use")

	pflag.String("token.backend", "jwt-rsa", "Token implementation to use")
	pflag.Duration("token.lifetime", time.Hour*10, "Token lifetime")
	pflag.Int("token.renewals", 5, "Maximum number of times the token may be renewed")

	pflag.Int("token.jwt.bits", 2048, "Bit length of generated keys")
	pflag.Bool("token.jwt.generate", false, "Generate keys if not available")

	pflag.Bool("pdb.watcher", false, "Enable the pdb filesystem watcher")
	pflag.Duration("pdb.watch-interval", 1 * time.Second, "Watch Interval")

	viper.SetDefault("server.port", 1729)
	viper.SetDefault("tls.certificate", "keys/tls.pem")
	viper.SetDefault("tls.key", "keys/tls.key")
}

func newServer() *rpc.NetAuthServer {
	// Need to setup the Database for use with the entity tree
	db, err := db.New()
	if err != nil {
		log.Fatalf("Fatal database error! (%s)", err)
	}
	log.Printf("Using %s", viper.GetString("db.backend"))

	crypto, err := crypto.New()
	if err != nil {
		log.Fatalf("Fatal crypto error! (%s)", err)
	}
	log.Printf("Using %s", viper.GetString("crypto.backend"))

	// Initialize the entity tree
	tree, err := tree.New(db, crypto)
	if err != nil {
		log.Fatalf("Fatal tree error! (%s)", err)
	}

	// Initialize the token service
	tokenService, err := token.New()
	if err != nil {
		log.Fatalf("Fatal error initializing token service: %s", err)
	}
	log.Printf("Using %s", viper.GetString("token.backend"))

	return &rpc.NetAuthServer{
		Tree:  tree,
		Token: tokenService,
	}
}

func loadConfig() {
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/netauth/")
	viper.AddConfigPath("$HOME/.netauth")
	viper.AddConfigPath(".")

	if *writeDefConfig != "" {
		if err := viper.WriteConfigAs(*writeDefConfig); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	// Attempt to load the config
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Fatal error reading configuration: ", err)
	}
}

func main() {
	// Do the config load before anything else, this might bail
	// out for a number of reasons.
	loadConfig()

	log.Println("NetAuth server is starting!")

	// Bind early so that if this fails we can just bail out.
	bindAddr := viper.GetString("server.bind")
	bindPort := viper.GetInt("server.port")
	sock, err := net.Listen("tcp", fmt.Sprintf("%s:%d", bindAddr, bindPort))
	if err != nil {
		log.Fatalf("Could not bind! %v", err)
	}
	log.Printf("Server bound on %s:%d", bindAddr, bindPort)

	// Setup the TLS parameters if necessary.
	var opts []grpc.ServerOption
	if !*insecure {
		cFile := viper.GetString("tls.certificate")
		ckFile := viper.GetString("tls.key")
		if !filepath.IsAbs(cFile) {
			cFile = filepath.Join(viper.GetString("core.home"), cFile)
		}
		if !filepath.IsAbs(ckFile) {
			ckFile = filepath.Join(viper.GetString("core.home"), ckFile)
		}
		log.Printf("TLS with the certificate %s and key %s", cFile, ckFile)
		creds, err := credentials.NewServerTLSFromFile(cFile, ckFile)
		if err != nil {
			log.Fatalf("TLS credentials could not be loaded! %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	} else {
		// Not using TLS in an auth server?  For shame...
		log.Println("===================================================================")
		log.Println("  WARNING WARNING WARNING WARNING WARNING WARNING WARNING WARNING  ")
		log.Println("===================================================================")
		log.Println("")
		log.Println("Launching without TLS! Your passwords will be shipped in the clear!")
		log.Println("Seriously, the option is --PWN_ME for a reason, you're trusting the")
		log.Println("network fabric with your authentication information, and this is a ")
		log.Println("bad idea.  Anyone on your local network can get passwords, tokens, ")
		log.Println("and other secure information.  You should instead obtain a ")
		log.Println("certificate and key and start the server with those.")
		log.Println("")
		log.Println("===================================================================")
		log.Println("  WARNING WARNING WARNING WARNING WARNING WARNING WARNING WARNING  ")
		log.Println("===================================================================")
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
		if !strings.Contains(*bootstrap, ":") {
			log.Fatal("Bootstrap string must be of the form <entity>:<secret>")
		}
		log.Println("Commencing Bootstrap")
		eParts := strings.Split(*bootstrap, ":")
		srv.Tree.Bootstrap(eParts[0], eParts[1])
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

	// Commence serving
	grpcServer.Serve(sock)
}
