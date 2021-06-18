package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	cfg string

	rootCmd = &cobra.Command{
		Use:   "nsutil <tool> [args]",
		Short: "Perform offline server actions.",
		Long:  rootCmdLongDocs,
	}

	rootCmdLongDocs = `
This utility can be used to perform certain tasks related to the
backing datastore for NetAuth.

DO NOT USE WHILE THE SERVER IS RUNNING.
`
)

func init() {
	viper.SetEnvPrefix("netauth")

	cobra.OnInitialize(onInit)
	rootCmd.PersistentFlags().StringVar(&cfg, "config", "", "Use an alternate config file")
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
}

func execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
