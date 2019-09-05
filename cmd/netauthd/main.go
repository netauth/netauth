package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/NetAuth/NetAuth/internal/crypto"
	_ "github.com/NetAuth/NetAuth/internal/crypto/all"
	"github.com/NetAuth/NetAuth/internal/db"
	_ "github.com/NetAuth/NetAuth/internal/db/all"
	plugin "github.com/NetAuth/NetAuth/internal/plugin/tree/manager"
	"github.com/NetAuth/NetAuth/internal/token"
	_ "github.com/NetAuth/NetAuth/internal/token/all"

	"github.com/NetAuth/NetAuth/internal/rpc"
	"github.com/NetAuth/NetAuth/internal/tree"
	_ "github.com/NetAuth/NetAuth/internal/tree/hooks"

	"github.com/hashicorp/go-hclog"
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

	version = "dev"
	commit  = "none"
	date    = "unknown"

	appLogger hclog.Logger
)

func init() {
	appLogger = hclog.New(&hclog.LoggerOptions{
		Name:  "netauthd",
		Level: hclog.LevelFromString("INFO"),
	})
	hclog.SetDefault(appLogger)

	pflag.String("tls.certificate", "keys/tls.crt", "Path to certificate file")
	pflag.String("tls.key", "keys/tls.key", "Path to key file")

	pflag.String("server.bind", "localhost", "Bind address, defaults to localhost")
	pflag.Int("server.port", 1729, "Serving port")
	pflag.String("core.home", "", "Base directory for NetAuth")

	pflag.String("db.backend", "ProtoDB", "Database implementation to use")

	pflag.String("crypto.backend", "bcrypt", "Cryptography system to use")

	pflag.String("token.backend", "jwt-rsa", "Token implementation to use")
	pflag.Duration("token.lifetime", time.Minute*10, "Token lifetime")

	pflag.Int("token.jwt.bits", 2048, "Bit length of generated keys")
	pflag.Bool("token.jwt.generate", false, "Generate keys if not available")

	pflag.Bool("pdb.watcher", false, "Enable the pdb filesystem watcher")
	pflag.Duration("pdb.watch-interval", 1*time.Second, "Watch Interval")

	pflag.String("log.level", "INFO", "Log verbosity level")

	viper.SetDefault("server.port", 1729)
	viper.SetDefault("tls.certificate", "keys/tls.pem")
	viper.SetDefault("tls.key", "keys/tls.key")
}

func newServer() *rpc.NetAuthServer {
	// Need to setup the Database for use with the entity tree
	db, err := db.New()
	if err != nil {
		appLogger.Error("Fatal database error", "error", err)
		os.Exit(1)
	}
	appLogger.Info("Database initialized", "backend", viper.GetString("db.backend"))

	crypto, err := crypto.New()
	if err != nil {
		appLogger.Error("Fatal crypto error", "error", err)
		os.Exit(1)
	}
	appLogger.Info("Cryptography system initialized", "backend", viper.GetString("crypto.backend"))

	// Initialize the entity tree
	tree, err := tree.New(db, crypto)
	if err != nil {
		appLogger.Error("Fatal initialization error", "error", err)
		os.Exit(1)
	}

	// Initialize the token service
	tokenService, err := token.New()
	if err != nil {
		appLogger.Error("Fatal token error", "error", err)
		os.Exit(1)
	}
	appLogger.Info("Token backend successfully initialized", "backend", viper.GetString("token.backend"))

	return &rpc.NetAuthServer{
		Tree:  tree,
		Token: tokenService,
		Log:   appLogger.Named("rpc"),
	}
}

func loadConfig() {
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.netauth")
	viper.AddConfigPath("/etc/netauth/")

	if *writeDefConfig != "" {
		if err := viper.WriteConfigAs(*writeDefConfig); err != nil {
			appLogger.Error("Error writing configuration", "error", err)
			os.Exit(2)
		}
		os.Exit(0)
	}

	// Attempt to load the config
	if err := viper.ReadInConfig(); err != nil {
		appLogger.Error("Fatal error reading configuration", "error", err)
		os.Exit(1)
	}
}

func main() {
	// Do the config load before anything else, this might bail
	// out for a number of reasons.
	loadConfig()
	appLogger.SetLevel(hclog.LevelFromString(viper.GetString("log.level")))

	appLogger.Info("NetAuth server is starting!")
	appLogger.Debug("Build information as follows", "version", version, "commit", commit, "builddate", date)

	// Bind early so that if this fails we can just bail out.
	bindAddr := viper.GetString("server.bind")
	bindPort := viper.GetInt("server.port")
	sock, err := net.Listen("tcp", fmt.Sprintf("%s:%d", bindAddr, bindPort))
	if err != nil {
		appLogger.Error("Could not bind!", "address", bindAddr, "port", bindPort)
		os.Exit(1)
	}
	appLogger.Debug("Server bind successful", "address", bindAddr, "port", bindPort)

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
		appLogger.Debug("TLS Enabled", "certificate", cFile, "key", ckFile)
		creds, err := credentials.NewServerTLSFromFile(cFile, ckFile)
		if err != nil {
			appLogger.Error("TLS could not be initialized", "error", err)
			os.Exit(1)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	} else {
		// Not using TLS in an auth server?  For shame...
		appLogger.Warn("===================================================================")
		appLogger.Warn("  WARNING WARNING WARNING WARNING WARNING WARNING WARNING WARNING  ")
		appLogger.Warn("===================================================================")
		appLogger.Warn("")
		appLogger.Warn("Launching without TLS! Your passwords will be shipped in the clear!")
		appLogger.Warn("Seriously, the option is --PWN_ME for a reason, you're trusting the")
		appLogger.Warn("network fabric with your authentication information, and this is a ")
		appLogger.Warn("bad idea.  Anyone on your local network can get passwords, tokens, ")
		appLogger.Warn("and other secure information.  You should instead obtain a ")
		appLogger.Warn("certificate and key and start the server with those.")
		appLogger.Warn("")
		appLogger.Warn("===================================================================")
		appLogger.Warn("  WARNING WARNING WARNING WARNING WARNING WARNING WARNING WARNING  ")
		appLogger.Warn("===================================================================")
	}

	// Spit out what backends we know about
	appLogger.Info("The following DB backends are registered:")
	for _, b := range db.GetBackendList() {
		appLogger.Info(fmt.Sprintf("  %s", b))
	}

	// Spit out what crypto backends we know about
	appLogger.Info("The following crypto implementations are registered:")
	for _, b := range crypto.GetBackendList() {
		appLogger.Info(fmt.Sprintf("  %s", b))
	}

	// Spit out the token services we know about
	appLogger.Info("The following token services are registered:")
	for _, b := range token.GetBackendList() {
		appLogger.Info(fmt.Sprintf("  %s", b))
	}

	// Get a plugin manager for extensibility
	p, err := plugin.New()
	if err != nil {
		appLogger.Warn("Problem initializing plugin manager", "error", err)
	}
	if viper.GetBool("plugin.enabled") {
		appLogger.Debug("Initializing tree plugins")
		p.LoadPlugins()
		p.RegisterEntityHooks()
		p.RegisterGroupHooks()
	} else {
		appLogger.Debug("Not running with plguins")
	}

	// Init the new server instance
	srv := newServer()

	if viper.GetBool("plugin.enabled") {
		p.ConfigureEntityChains(srv.Tree.RegisterEntityHookToChain)
		p.ConfigureGroupChains(srv.Tree.RegisterGroupHookToChain)
	}

	// Attempt to bootstrap a superuser
	if len(*bootstrap) != 0 {
		if !strings.Contains(*bootstrap, ":") {
			appLogger.Error("Bootstrap string must be in the format of <entity>:<secret>")
			os.Exit(1)
		}
		appLogger.Info("Beginning Bootstrap")
		eParts := strings.Split(*bootstrap, ":")
		srv.Tree.Bootstrap(eParts[0], eParts[1])
		appLogger.Info("Bootstrap complete")
	}

	// If it wasn't used make sure its disabled since it can
	// create arbitrary root users.
	srv.Tree.DisableBootstrap()

	// Instantiate and launch.  This will block and the server
	// will serve forever.
	appLogger.Info("Ready to Serve...")
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterNetAuthServer(grpcServer, srv)

	// Register shutdown machinery
	c := make(chan os.Signal, 1)
	go func() {
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	}()

	done := make(chan struct{}, 1)
	go func() {
		<-c
		appLogger.Info("Shutting down...")
		grpcServer.GracefulStop()
		p.Shutdown()
		close(done)
	}()

	// Commence serving
	grpcServer.Serve(sock)

	appLogger.Info("Waiting for shutdown to complete")
	<-done
	appLogger.Info("Goodbye!")
}
