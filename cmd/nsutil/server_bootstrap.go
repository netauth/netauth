package main

import (
	"context"
	"fmt"
	"os"

	"github.com/bgentry/speakeasy"
	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/netauth/netauth/internal/crypto"
	_ "github.com/netauth/netauth/internal/crypto/bcrypt"
	"github.com/netauth/netauth/internal/db"
	_ "github.com/netauth/netauth/internal/db/bitcask"
	_ "github.com/netauth/netauth/internal/db/filesystem"
	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"
	_ "github.com/netauth/netauth/internal/tree/hooks"

	pb "github.com/netauth/protocol"
)

var (
	serverBootstrapCmd = &cobra.Command{
		Use:   "bootstrap <username>",
		Short: "Make or update a user and provide them root authority",
		Long:  serverBootstrapCmdLongDocs,
		Run:   serverBootstrapCmdRun,
		Args:  cobra.ExactArgs(1),
	}

	serverBootstrapCmdLongDocs = `
The bootstrap command will either create the specified user or update
an existing one of the same ID.  The password will be reset and the
user will be directly assigned the GLOBAL_ROOT capability flag, which
will permit further bootstrapping tasks.

!!! ACHTUNG !!!
You must only run this command with the server stopped to ensure your
data storage remains consistent.
`
)

func init() {
	rootCmd.AddCommand(serverBootstrapCmd)
}

func serverBootstrapCmdRun(c *cobra.Command, args []string) {
	hclog.L().SetLevel(hclog.LevelFromString("TRACE"))
	crypto.SetParentLogger(hclog.L())
	db.SetParentLogger(hclog.L())
	tree.SetParentLogger(hclog.L())
	startup.DoCallbacks()
	ctx := context.Background()

	dbImpl, err := db.New(viper.GetString("db.backend"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal database error: %s\n", err)
		os.Exit(1)
	}
	cryptoImpl, err := crypto.New(viper.GetString("crypto.backend"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal crypto error: %s\n", err)
		os.Exit(1)
	}
	tree, err := tree.New(dbImpl, cryptoImpl, hclog.L())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal initialization error: %s\n", err)
		os.Exit(1)
	}

	_, err = tree.FetchEntity(ctx, args[0])
	if err != db.ErrUnknownEntity && err != nil {
		fmt.Fprintf(os.Stderr, "Error checking for entity existence: %s\n", err)
	}
	if err == db.ErrUnknownEntity {
		if err := tree.CreateEntity(ctx, args[0], -1, ""); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating entity: %s\n", err)
			os.Exit(1)
		}
	}

	// Reset the password
	secret, err := speakeasy.Ask(fmt.Sprintf("New secret for %s: ", args[0]))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error prompting secret: %s", err)
		os.Exit(1)
	}
	if err := tree.SetSecret(ctx, args[0], secret); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting secret: %s\n", err)
		os.Exit(1)
	}

	// Ensure the entity is unlocked
	if err := tree.UnlockEntity(ctx, args[0]); err != nil {
		fmt.Fprintf(os.Stderr, "Error unlocking entity: %s\n", err)
		os.Exit(1)
	}

	// Bestow GLOBAL_ROOT capability
	if err := tree.SetEntityCapability2(ctx, args[0], pb.Capability_GLOBAL_ROOT.Enum()); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting GLOBAL_ROOT: %s\n", err)
	}
}
