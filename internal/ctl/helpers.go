package ctl

import (
	"fmt"
	"strings"

	"github.com/bgentry/speakeasy"
	"github.com/spf13/viper"

	"github.com/NetAuth/NetAuth/pkg/client"

	pb "github.com/NetAuth/Protocol"
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

func getToken(c *client.NetAuthClient, entity string) (string, error) {
	t, err := c.GetToken(entity, "")
	switch err {
	case nil:
		return t, nil
	case client.ErrTokenUnavailable:
		return c.GetToken(entity, getSecret(""))
	default:
		return "", err
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

	startChar := ""
	for _, f := range fieldList {
		switch strings.ToLower(f) {
		case "id":
			fmt.Printf("%sID: %s\n", startChar, entity.GetID())
		case "number":
			fmt.Printf("%sNumber: %d\n", startChar, entity.GetNumber())
		case "primarygroup":
			if entity.Meta != nil && entity.GetMeta().GetPrimaryGroup() != "" {
				fmt.Printf("%sPrimary Group: %s\n", startChar, entity.GetMeta().GetPrimaryGroup())
			}
		case "gecos":
			if entity.Meta != nil && entity.GetMeta().GetGECOS() != "" {
				fmt.Printf("%sGECOS: %s\n", startChar, entity.GetMeta().GetGECOS())
			}
		case "legalname":
			if entity.Meta != nil && entity.GetMeta().GetLegalName() != "" {
				fmt.Printf("%slegalName: %s\n", startChar, entity.GetMeta().GetLegalName())
			}
		case "displayname":
			if entity.Meta != nil && entity.Meta.GetDisplayName() != "" {
				fmt.Printf("%sdisplayname: %s\n", startChar, entity.GetMeta().GetDisplayName())
			}
		case "homedir":
			if entity.Meta != nil && entity.GetMeta().GetHome() != "" {
				fmt.Printf("%shomedir: %s\n", startChar, entity.GetMeta().GetHome())
			}
		case "shell":
			if entity.Meta != nil && entity.GetMeta().GetShell() != "" {
				fmt.Printf("%sshell: %s\n", startChar, entity.GetMeta().GetShell())
			}
		case "graphicalshell":
			if entity.Meta != nil && entity.GetMeta().GetGraphicalShell() != "" {
				fmt.Printf("%sgraphicalShell: %s\n", startChar, entity.GetMeta().GetGraphicalShell())
			}
		case "badgenumber":
			if entity.Meta != nil && entity.GetMeta().GetBadgeNumber() != "" {
				fmt.Printf("%sbadgeNumber: %s\n", startChar, entity.GetMeta().GetBadgeNumber())
			}
		case "capabilities":
			if entity.Meta != nil && len(entity.GetMeta().GetCapabilities()) != 0 {
				fmt.Printf("%sCapabilities (Direct):\n", startChar)
				for i := range entity.GetMeta().GetCapabilities() {
					fmt.Printf("%s  - %s\n", startChar, entity.GetMeta().GetCapabilities()[i])
				}
			}
		}
		startChar = "  "
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
			"expansions",
			"capabilities",
		}
	}

	startChar := ""
	for _, f := range fieldList {
		switch strings.ToLower(f) {
		case "name":
			fmt.Printf("%sName: %s\n", startChar, group.GetName())
		case "displayname":
			fmt.Printf("%sDisplay Name: %s\n", startChar, group.GetDisplayName())
		case "number":
			fmt.Printf("%sNumber: %d\n", startChar, group.GetNumber())
		case "managedby":
			if group.GetManagedBy() == "" {
				continue
			}
			fmt.Printf("%sManaged By: %s\n", startChar, group.GetManagedBy())
		case "expansions":
			for _, exp := range group.GetExpansions() {
				fmt.Printf("%sExpansion: %s\n", startChar, exp)
			}
		case "capabilities":
			if len(group.GetCapabilities()) != 0 {
				fmt.Printf("%sCapabilities:\n", startChar)
				for i := range group.GetCapabilities() {
					fmt.Printf("%s  - %s\n", startChar, group.GetCapabilities()[i])
				}
			}
		}
		startChar = "  "
	}
}
