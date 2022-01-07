package ctl

import (
	"context"
	"fmt"
	"os"
	"os/user"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/netauth/netauth/pkg/netauth"
	tkn "github.com/netauth/netauth/pkg/token"
	"github.com/netauth/netauth/pkg/token/cache"
	"github.com/netauth/netauth/pkg/token/keyprovider"
)

var (
	rpc *netauth.Client

	tcache cache.TokenCache
	tsvc   tkn.Service

	cfg        string
	rootEntity string
	secret     string

	ctx context.Context

	rootCmd = &cobra.Command{
		Use:              "netauth <subsystem> <command> [flags] [args]",
		Short:            "Interact with the NetAuth system.",
		Long:             rootCmdLongDocs,
		PersistentPreRun: initialize,
	}

	rootCmdLongDocs = `
NetAuth is an authentication and authorization system for small to
medium scale networks.  This tool is designed to be the root point of
interaction with the NetAuth system and is divided up into subsystems
and subcommands for interaction with specific facets of the NetAuth
ecosystem.`
)

func init() {
	viper.SetEnvPrefix("netauth")

	cobra.OnInitialize(onInit)
	rootCmd.PersistentFlags().StringVar(&cfg, "config", "", "Use an alternate config file")
	rootCmd.PersistentFlags().StringVar(&rootEntity, "entity", "", "Specify a non-default entity to make requests as")
	rootCmd.PersistentFlags().StringVar(&secret, "secret", "", "Specify the request secret on the command line")

	viper.BindPFlag("entity", rootCmd.PersistentFlags().Lookup("entity"))
	viper.BindEnv("entity")
	viper.BindPFlag("secret", rootCmd.PersistentFlags().Lookup("secret"))
	viper.BindEnv("secret")
}

func onInit() {
	viper.BindPFlags(pflag.CommandLine)
	if cfg != "" {
		viper.SetConfigFile(cfg)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.netauth")
		viper.AddConfigPath("/etc/netauth/")
	}
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config:", err)
		os.Exit(1)
	}
	viper.Set("client.ServiceName", "netauth")

	user, err := user.Current()
	if err != nil {
		fmt.Println("Could not get default user:", err)
	}
	viper.SetDefault("entity", user.Username)
}

// Execute serves as the entrypoint to the ctl package.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func initialize(*cobra.Command, []string) {
	var err error
	tcache, err = cache.NewTokenCache(viper.GetString("token.cache"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing token cache: %v", err)
		os.Exit(1)
	}

	kp, err := keyprovider.New(viper.GetString("token.keyprovider"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing token service: %v", err)
		os.Exit(1)
	}

	tsvc, err = tkn.New(viper.GetString("token.backend"), kp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing token service: %v", err)
		os.Exit(1)
	}

	rpc, err = netauth.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during client initialization: %v", err)
		os.Exit(1)
	}

	ctx = context.Background()
	rpc.SetServiceName("netauth")
}
