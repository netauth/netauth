package ctl

import (
	"fmt"
	"strings"

	pb "github.com/NetAuth/NetAuth/pkg/proto"
)

// ensureSecret prompts for the secret if it was not provided already.
// This gets around the secret being visible on the command line.
func ensureSecret() {
	if secret == "" {
		fmt.Print("Secret: ")
		_, err := fmt.Scanln(&secret)
		if err != nil {
			fmt.Printf("Error: %s", err)
		}
	}
}


func printEntity(entity *pb.Entity, fields string) {
	fieldList := []string{}
	if fields != "" {
		fieldList = strings.Split(fields, ",")
	} else {
		fieldList = []string{
			"ID",
			"uidNumber",
			"GECOS",
			"legalName",
			"displayName",
			"homedir",
			"shell",
			"graphicalShell",
			"badgeNumber",
		}
	}

	for _, f := range fieldList {
		switch f {
		case "ID":
			fmt.Printf("ID: %s\n", entity.GetID())
		case "uidNumber":
			fmt.Printf("uidNumber: %d\n", entity.GetUidNumber())
		case "GECOS":
			if entity.Meta != nil && entity.GetMeta().GetGECOS() != "" {
				fmt.Printf("GECOS: %s\n", entity.GetMeta().GetGECOS())
			}
		case "legalName":
			if entity.Meta != nil && entity.GetMeta().GetLegalName() != "" {
				fmt.Printf("legalName: %s\n", entity.GetMeta().GetLegalName())
			}
		case "displayName":
			if entity.Meta != nil && entity.Meta.GetDisplayName() != "" {
				fmt.Printf("displayname: %s\n", entity.GetMeta().GetDisplayName())
			}
		case "homedir":
			if  entity.Meta != nil && entity.GetMeta().GetHomedir() != "" {
				fmt.Printf("homedir: %s\n", entity.GetMeta().GetHomedir())
			}
		case "shell":
			if entity.Meta != nil && entity.GetMeta().GetShell() != "" {
				fmt.Printf("shell: %s\n", entity.GetMeta().GetShell())
			}
		case "graphicalShell":
			if entity.Meta != nil && entity.GetMeta().GetGraphicalShell() != "" {
				fmt.Printf("graphicalShell: %s\n", entity.GetMeta().GetGraphicalShell())
			}
		case "badgeNumber":
			if entity.Meta != nil && entity.GetMeta().GetBadgeNumber() != "" {
				fmt.Printf("badgeNumber: %s\n", entity.GetMeta().GetBadgeNumber())
			}
		}
	}

}
