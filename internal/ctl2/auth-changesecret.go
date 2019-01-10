package ctl2

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/NetAuth/NetAuth/pkg/client"
)

var (
	csEntity string
	csSecret string

	authChangeSecretCmd = &cobra.Command{
		Use:     "change-secret",
		Short:   "Change an entity secret",
		Example: authChangeSecretExample,
		Long:    authChangeSecretLongDocs,
		Run:     authChangeSecretRun,
	}

	authChangeSecretLongDocs = `
The change-secret command is used to change an entity's secret either
reflexively (the entity requests the change) or administratively
(another entity changes the secret).`

	authChangeSecretExample = `$ netauth auth change-secret
Old Secret:
New Secret:
Verify Secret:
Secret Changed

$ netauth auth change-secret --csEntity demo
New Secret:
Verify Secret:
Secret Changed`
)

func init() {
	authCmd.AddCommand(authChangeSecretCmd)
	authChangeSecretCmd.Flags().StringVar(&csEntity, "csEntity", "", "Entity to change secret")
	authChangeSecretCmd.Flags().StringVar(&csSecret, "csSecret", "", "Secret (omit for prompt)")
}

func authChangeSecretRun(cmd *cobra.Command, args []string) {
	s := ""
	t := ""

	// Grab a client
	c, err := client.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Self change if unset
	if csEntity == "" {
		csEntity = viper.GetString("entity")
	}

	// Get either secret or token
	if csEntity == viper.GetString("entity") {
		s = getSecret("Old Secret: ")
	} else {
		// Get the authorization token
		t, err = getToken(c, viper.GetString("entity"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// Get the secret if it wasn't specified on the line
	if csSecret == "" {
		one := getSecret("New Secret: ")
		two := getSecret("Verify Secret: ")

		if one != two {
			fmt.Println("Secrets do not match!")
			os.Exit(1)
		}

		csSecret = one
	}

	// Change the secret
	result, err := c.ChangeSecret(viper.GetString("entity"), s, csEntity, csSecret, t)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(result.GetMsg())
}
