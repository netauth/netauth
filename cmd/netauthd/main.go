package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/netauth/netauth/internal/crypto"
	_ "github.com/netauth/netauth/internal/crypto/all"
	"github.com/netauth/netauth/internal/db"
	_ "github.com/netauth/netauth/internal/db/all"
	plugin "github.com/netauth/netauth/internal/plugin/tree/manager"
	"github.com/netauth/netauth/internal/token"
	_ "github.com/netauth/netauth/internal/token/all"

	"github.com/netauth/netauth/internal/rpc"
	"github.com/netauth/netauth/internal/rpc2"
	"github.com/netauth/netauth/internal/tree"
	_ "github.com/netauth/netauth/internal/tree/hooks"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/NetAuth/Protocol"
	rpb "github.com/NetAuth/Protocol/v2"
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
	viper.SetDefault("plugin.path", filepath.Join(viper.GetString("core.home"), "plugins"))
}

// newServer sets up a version 1 RPC server interface.  This should
// really be a part of internal/rpc, but that entire interface is
// deprecated and so no refactoring will happen to it.
func newServer(tree rpc.EntityTree, tokenService token.Service) *rpc.NetAuthServer {
	return &rpc.NetAuthServer{
		Tree:  tree,
		Token: tokenService,
		Log:   appLogger.Named("rpc"),
	}
}

// newSocket binds the listening socket to the ports specified in the
// configuration file.
func newSocket() (net.Listener, error) {
	// Bind early so that if this fails we can just bail out.
	bindAddr := viper.GetString("server.bind")
	bindPort := viper.GetInt("server.port")
	sock, err := net.Listen("tcp", fmt.Sprintf("%s:%d", bindAddr, bindPort))
	if err != nil {
		appLogger.Error("Could not bind!", "address", bindAddr, "port", bindPort)
		return nil, err
	}
	appLogger.Debug("Server bind successful", "address", bindAddr, "port", bindPort)
	return sock, nil
}

// newGRPCServer takes care of setting up a grpc.Server to bind
// implementations into.  This includes loading certificate files if
// serving with TLS, or printing a large scary warning if transport
// security has been intentionally disabled.
func newGRPCServer() (*grpc.Server, error) {
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
			return nil, err
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
	grpcServer := grpc.NewServer(opts...)
	return grpcServer, nil
}

// loadConfig is a convenience function that handles the loading of
// the viper configuration singleton.  This function is called just
// after flag parsing completes and if it is unsuccessful it aborts
// the entire server process.
func loadConfig() error {
	viper.BindPFlags(pflag.CommandLine)
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.netauth")
	viper.AddConfigPath("/etc/netauth/")

	// Attempt to load the config
	if err := viper.ReadInConfig(); err != nil {
		appLogger.Error("Fatal error reading configuration", "error", err)
		return err
	}
	return nil
}

// writeDefaultconfig spits out the default config to the requested
// file.  This is used instead of just calling an unconditional
// SafeWriteConfigAs because it is assumed that the user that runs the
// server won't have write permissions to the real configuration file.
// Using this argument allows the operator to explicitly state where
// they wish the config example to be written.
func writeDefaultConfig() error {
	if err := viper.WriteConfigAs(*writeDefConfig); err != nil {
		appLogger.Error("Error writing configuration", "error", err)
		return err
	}
	return nil
}

// doInfoLogScroll prints out the information about the binary version
// and what statically compiled implementations are available.
func doInfoLogScroll() {
	appLogger.Info("NetAuth server is starting!")
	appLogger.Debug("Server home directory", "directory", viper.GetString("core.home"))
	appLogger.Debug("Build information as follows", "version", version, "commit", commit, "builddate", date)

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
}

// doPluginEarlySetup takes care of some setup tasks needed to create
// a plugin manager for tree plugins.  This function will only perform
// these setup tasks if plugins are enabled, otherwise it will not
// alter the configuration of the server.
func doPluginEarlySetup() plugin.Manager {
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
	return p
}

// doSuperuserBootstrap handles the text parsing to create a new
// superuser in the system.  Most importantly this function ensures
// that the format is correct for the bootstrap string prior to trying
// to split it.
func doSuperuserBootstrap(t *tree.Manager) error {
	if !strings.Contains(*bootstrap, ":") {
		return errors.New("bootstrap string must be in the format of <entity>:<secret>")
	}
	appLogger.Info("Beginning Bootstrap")
	eParts := strings.SplitN(*bootstrap, ":", 2)
	t.Bootstrap(eParts[0], eParts[1])
	appLogger.Info("Bootstrap complete")
	return nil
}

func main() {
	// Parse flags first, this is required to be able to chose
	// whether or not to write out the default configuration
	// rather than starting the server.
	pflag.Parse()

	// The configuration format is fairly well documented, but it
	// is also useful to be able to write out the configuration
	// file with all the default options.  When the config file is
	// being written the server will not start.
	if *writeDefConfig != "" {
		writeDefaultConfig()
		os.Exit(0)
	}

	// Load config as early as possible.  Some of the lower
	// initialization sections expect this to just work, so if any
	// errors are encountered we take the whole server down here.
	if err := loadConfig(); err != nil {
		os.Exit(1)
	}
	appLogger.SetLevel(hclog.LevelFromString(viper.GetString("log.level")))

	// This spits out all the bootup information, debugging
	// tokens, and some other diagnostic information that make up
	// the first 40 or so lines of the startup log.
	doInfoLogScroll()

	// The plugin system requires initialization very early in the
	// server startup.  This will scan for and register external
	// plugins, and register the hooks that support the plugin
	// system.
	pluginManager := doPluginEarlySetup()

	// The data storage layer and cryptographic engine are next to
	// initialize.  These modules provide core services to the
	// entity tree which initializes immediately afterwards.
	dbImpl, err := db.New()
	if err != nil {
		appLogger.Error("Fatal database error", "error", err)
		os.Exit(1)
	}
	appLogger.Info("Database initialized", "backend", viper.GetString("db.backend"))

	cryptoImpl, err := crypto.New()
	if err != nil {
		appLogger.Error("Fatal crypto error", "error", err)
		os.Exit(1)
	}
	appLogger.Info("Cryptography system initialized", "backend", viper.GetString("crypto.backend"))

	// The Tree is the core component of the server.  Its the part
	// that actually provides the interface for working with
	// entities, working with groups, and defining the
	// relationships between the two.  If the plugin system is
	// being used, then the tree action configurations (chains)
	// need to be reconfigured to enable the external plugin
	// hooks.
	tree, err := tree.New(dbImpl, cryptoImpl)
	if err != nil {
		appLogger.Error("Fatal initialization error", "error", err)
		os.Exit(1)
	}
	if viper.GetBool("plugin.enabled") {
		pluginManager.ConfigureEntityChains(tree.RegisterEntityHookToChain)
		pluginManager.ConfigureGroupChains(tree.RegisterGroupHookToChain)
	}

	// Bootstrapping a new server, or recovering from certain
	// scenarios requires being able to generate a new user, or
	// elevate an existing one, to the global superuser state.
	// This action requires the user to have the root on the
	// NetAuth server, and this is generally accepted to be enough
	// to override the internal permissions model based on
	// external authority.  After this action is taken or not
	// taken, the capability to bootstrap is disabled for the
	// lifetime of the server.
	if len(*bootstrap) != 0 {
		if err := doSuperuserBootstrap(tree); err != nil {
			appLogger.Error("Critical error during superuser bootstrap", "error", err)
			os.Exit(1)
		}
	}
	tree.DisableBootstrap()

	// NetAuth's internal security model is token based.  The
	// token service is distinct from the tree, and can wait to
	// come online until the tree has been initiailized (and by
	// extension the plugin system).
	tokenService, err := token.New()
	if err != nil {
		appLogger.Error("Fatal token error", "error", err)
		os.Exit(1)
	}
	appLogger.Info("Token backend successfully initialized", "backend", viper.GetString("token.backend"))

	// Initializing the gRPC Server happens only once the
	// primitives that it will consume have been initialized.  At
	// the point that the gRPC components initialize, TLS keys
	// will be loaded.  If the server is being run in an insecure
	// mode then a warning will be printed to the log before an
	// insecure server is returned.
	grpcServer, err := newGRPCServer()
	if err != nil {
		os.Exit(1)
	}

	// A NetAuth server may serve more than one protocol version
	// at a time.  This section binds the different application
	// protocol versions to the grpcServer.
	pb.RegisterNetAuthServer(
		grpcServer,
		newServer(
			tree,
			tokenService,
		),
	)
	rpb.RegisterNetAuth2Server(
		grpcServer,
		rpc2.New(
			rpc2.Refs{
				TokenService: tokenService,
				Tree:         tree,
			},
		),
	)

	// While the server is for the most part stateless, the
	// plugins might not be.  This block registers the shutdown
	// machinery that allows the server to make a clean exit and
	// not leak processes.  Any additional parallel shutdown tasks
	// should be added to the goroutine which will be signalled in
	// the event of a process interrupt or termination signal.
	c := make(chan os.Signal, 1)
	go func() {
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	}()
	done := make(chan struct{}, 1)
	go func() {
		<-c
		appLogger.Info("Shutting down...")
		grpcServer.GracefulStop()
		pluginManager.Shutdown()
		close(done)
	}()

	// Commence serving.  This call is blocking and is only
	// interrupted by the shutdown call being made above which
	// will only happen if an external process supervisor signals
	// the server to shutdown.  While it might seem odd to bind
	// the server here, its comparatively more likely that file
	// permissions will be wrong on some important file earlier on
	// than the port won't bind.
	appLogger.Info("Ready to Serve...")
	sock, err := newSocket()
	if err != nil {
		os.Exit(1)
	}
	grpcServer.Serve(sock)

	// Once the server has been signalled to shut down, it is
	// necessary to wait for the parallel shutdown tasks to signal
	// they are done as well.  These tasks handle the closing down
	// of supervised processes and other critical components that
	// could leak to the operating system if not shut down
	// correctly, so we'll wait patiently here for these tasks to
	// complete.
	appLogger.Info("Waiting for shutdown to complete")
	<-done
	appLogger.Info("Goodbye!")
}
