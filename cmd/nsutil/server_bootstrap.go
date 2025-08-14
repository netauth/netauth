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
	noRebootstrap      bool
	serverBootstrapCmd = &cobra.Command{
		Use:   "bootstrap [--no-rebootstrap] <username>",
		Short: "Make or update a user and provide them root authority",
		Long:  serverBootstrapCmdLongDocs,
		Run:   serverBootstrapCmdRun,
		Args: func(cmd *cobra.Command, args []string) error {
			switch len(args) {
			case 0:
				if os.Getenv("NETAUTH_UNATTENDED_BOOTSTRAP_NAME") == "" {
					return fmt.Errorf("Need an entity name either on the commandline or the NETAUTH_UNATTENDED_BOOTSTRAP_NAME environment variable")
				}
				return nil
			case 1:
				if args[0] == "--no-rebootstrap" && os.Getenv("NETAUTH_UNATTENDED_BOOTSTRAP_NAME") == "" {
					return fmt.Errorf("Need an entity name either on the commandline or the NETAUTH_UNATTENDED_BOOTSTRAP_NAME environment variable")
				}
				return nil
			case 2:
				if args[0] == "--no-rebootstrap" || args[1] == "--no-rebootstrap" {
					return nil
				}
				return fmt.Errorf("The only valid arguments are the entity name and the --no-rebootstrap flag")
			default:
				return fmt.Errorf("The only valid flag is --no-rebootstrap and the user is optional")
			}
		},
	}

	serverBootstrapCmdLongDocs = `
The bootstrap command will either create the specified user or update
an existing one of the same ID.  The password will be reset and the
user will be directly assigned the GLOBAL_ROOT capability flag, which
will permit further bootstrapping tasks.

Bootstrap will optionally read the bootstrap secret and user from
environment variables:
NETAUTH_UNATTENDED_BOOTSTRAP_NAME
NETAUTH_UNATTENDED_BOOTSTRAP_SECRET

!!! ACHTUNG !!!
You must only run this command with the server stopped to ensure your
data storage remains consistent.
`
)

func init() {
	serverBootstrapCmd.Flags().BoolVarP(&noRebootstrap, "no-rebootstrap", "", false, "Skip entity unlocking if it already exists")
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

	opts := []tree.Option{
		tree.WithStorage(dbImpl),
		tree.WithCrypto(cryptoImpl),
		tree.WithLogger(hclog.L()),
	}

	tree, err := tree.New(opts...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal initialization error: %s\n", err)
		os.Exit(1)
	}

	entity := os.Getenv("NETAUTH_UNATTENDED_BOOTSTRAP_NAME")
	if entity == "" {
		entity = args[0]
	}

	_, err = tree.FetchEntity(ctx, entity)
	if err != db.ErrUnknownEntity && err != nil {
		fmt.Fprintf(os.Stderr, "Error checking for entity existence: %s\n", err)
	}
	if err == db.ErrUnknownEntity {
		if err := tree.CreateEntity(ctx, entity, -1, ""); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating entity: %s\n", err)
			os.Exit(1)
		}
	} else if noRebootstrap {
		appLogger.Info("Bootstrap killed early since the entity exists")
		os.Exit(0)
	}

	// Reset the password
	secret := os.Getenv("NETAUTH_UNATTENDED_BOOTSTRAP_SECRET")
	if secret == "" {
		secret, err = speakeasy.Ask(fmt.Sprintf("New secret for %s: ", entity))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error prompting secret: %s", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Secret for %s provided via the environment", entity)
	}
	if err := tree.SetSecret(ctx, entity, secret); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting secret: %s\n", err)
		os.Exit(1)
	}

	// Ensure the entity is unlocked
	if err := tree.UnlockEntity(ctx, entity); err != nil {
		fmt.Fprintf(os.Stderr, "Error unlocking entity: %s\n", err)
		os.Exit(1)
	}

	// Bestow GLOBAL_ROOT capability
	if err := tree.SetEntityCapability2(ctx, entity, pb.Capability_GLOBAL_ROOT.Enum()); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting GLOBAL_ROOT: %s\n", err)
	}

	appLogger.Info("Bootstrap complete, enter when ready")
}
