package ctl

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	cfg        string
	rootEntity string
	secret     string

	rootCmd = &cobra.Command{
		Use:   "netauth <subsystem> <command> [flags] [args]",
		Short: "Interact with the NetAuth system.",
		Long:  rootCmdLongDocs,
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
	if _, set := os.LookupEnv("NETAUTH_VERBOSE"); !set {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
