package ctl

import (
	"fmt"
	"os"
	"strings"

	"github.com/bgentry/speakeasy"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/netauth/netauth/pkg/token/cache"

	pb "github.com/netauth/protocol"
)

// Prompt for the secret if it wasn't provided in cleartext.
func getSecret(prompt string) string {
	if prompt == "" {
		prompt = "Secret: "
	}

	if viper.GetString("secret") != "" {
		return viper.GetString("secret")
	}
	secret, err := speakeasy.Ask(prompt)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	return secret
}

// token is used exclusively by the CLI to provide tokens either from
// the cache or the RPC call. It returns a string or calls exit, there
// are no conditions where the string will be returned without a
// token.  Since this is for the CLI, it always uses the value of the
// entity that the call is being made as.
func token() string {
	t, err := tcache.GetToken(viper.GetString("entity"))
	switch {
	case err == cache.ErrNoCachedToken:
		return refreshToken()
	case tokenIsExpired(t):
		return refreshToken()
	case !tokenIsExpired(t):
		return t
	default:
		return ""
	}
}

// tokenIsExpired checks if a token is no longer valid.  Technically
// it checks if the CLI can validate it, but in this case we can treat
// a validation error as cause to renew it.
func tokenIsExpired(t string) bool {
	return rpc.AuthValidateToken(ctx, t) != nil
}

// refreshToken is a convenience function to acquire a token or die
// trying.  This is meant for CLI use only, and thus we call exit here
// if necessary to handle errors.
func refreshToken() string {
	return refreshTokenWithSecret(getSecret(""))
}

// refreshTokenWithSecret performs an immediate refresh of the token.
func refreshTokenWithSecret(secret string) string {
	t, err := rpc.AuthGetToken(ctx, viper.GetString("entity"), secret)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := tcache.PutToken(viper.GetString("entity"), t); err != nil {
		fmt.Fprintf(os.Stderr, "Error caching token: %v\n", err)
	}
	return t
}

func kvArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("this command takes at least 3 arguments")
	}
	action := strings.ToUpper(args[1])
	if action == "UPSERT" && len(args) != 4 {
		return fmt.Errorf("upsert requires a key and a value")
	}

	switch action {
	case "UPSERT":
		return nil
	case "CLEARFUZZY":
		return nil
	case "CLEAREXACT":
		return nil
	case "READ":
		return nil
	default:
		return fmt.Errorf("action must be one of UPSERT, CLEARFUZZY, CLEAREXACT, or READ")
	}
}

func printEntity(entity *pb.Entity, fields string) {
	var fieldList []string

	if fields != "" {
		fieldList = strings.Split(fields, ",")
	} else {
		fieldList = []string{
			"ID",
			"number",
			"PrimaryGroup",
			"GECOS",
			"legalName",
			"displayName",
			"homedir",
			"shell",
			"graphicalShell",
			"badgeNumber",
			"capabilities",
		}
	}

	for _, f := range fieldList {
		switch strings.ToLower(f) {
		case "id":
			fmt.Printf("ID: %s\n", entity.GetID())
		case "number":
			fmt.Printf("Number: %d\n", entity.GetNumber())
		case "primarygroup":
			if entity.Meta != nil && entity.GetMeta().GetPrimaryGroup() != "" {
				fmt.Printf("Primary Group: %s\n", entity.GetMeta().GetPrimaryGroup())
			}
		case "gecos":
			if entity.Meta != nil && entity.GetMeta().GetGECOS() != "" {
				fmt.Printf("GECOS: %s\n", entity.GetMeta().GetGECOS())
			}
		case "legalname":
			if entity.Meta != nil && entity.GetMeta().GetLegalName() != "" {
				fmt.Printf("legalName: %s\n", entity.GetMeta().GetLegalName())
			}
		case "displayname":
			if entity.Meta != nil && entity.Meta.GetDisplayName() != "" {
				fmt.Printf("displayname: %s\n", entity.GetMeta().GetDisplayName())
			}
		case "homedir":
			if entity.Meta != nil && entity.GetMeta().GetHome() != "" {
				fmt.Printf("homedir: %s\n", entity.GetMeta().GetHome())
			}
		case "shell":
			if entity.Meta != nil && entity.GetMeta().GetShell() != "" {
				fmt.Printf("shell: %s\n", entity.GetMeta().GetShell())
			}
		case "graphicalshell":
			if entity.Meta != nil && entity.GetMeta().GetGraphicalShell() != "" {
				fmt.Printf("graphicalShell: %s\n", entity.GetMeta().GetGraphicalShell())
			}
		case "badgenumber":
			if entity.Meta != nil && entity.GetMeta().GetBadgeNumber() != "" {
				fmt.Printf("badgeNumber: %s\n", entity.GetMeta().GetBadgeNumber())
			}
		case "capabilities":
			if entity.Meta != nil && len(entity.GetMeta().GetCapabilities()) != 0 {
				fmt.Printf("Capabilities (Direct):\n")
				for i := range entity.GetMeta().GetCapabilities() {
					fmt.Printf("  - %s\n", entity.GetMeta().GetCapabilities()[i].String())
				}
			}
		}
	}
}

func printGroup(group *pb.Group, fields string) {
	var fieldList []string

	if fields != "" {
		fieldList = strings.Split(fields, ",")
	} else {
		fieldList = []string{
			"name",
			"displayName",
			"number",
			"managedBy",
			"rules",
			"capabilities",
		}
	}

	for _, f := range fieldList {
		switch strings.ToLower(f) {
		case "name":
			fmt.Printf("Name: %s\n", group.GetName())
		case "displayname":
			fmt.Printf("Display Name: %s\n", group.GetDisplayName())
		case "number":
			fmt.Printf("Number: %d\n", group.GetNumber())
		case "managedby":
			if group.GetManagedBy() == "" {
				continue
			}
			fmt.Printf("Managed By: %s\n", group.GetManagedBy())
		case "rules":
			for _, exp := range group.GetExpansions() {
				fmt.Printf("Rule: %s\n", exp)
			}
		case "capabilities":
			if len(group.GetCapabilities()) != 0 {
				fmt.Printf("Capabilities:\n")
				for i := range group.GetCapabilities() {
					fmt.Printf("  - %s\n", group.GetCapabilities()[i])
				}
			}
		}
	}
}
