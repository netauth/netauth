package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	keygenRSACmd = &cobra.Command{
		Use:   "rsa <bits>",
		Short: "Create an RSA Keypair",
		Long:  keygenRSACmdLongDocs,
		Run:   keygenRSACmdRun,
		Args:  cobra.ExactArgs(1),
	}

	keygenRSACmdLongDocs = `The RSA Keygen command allows you to generate a keypair suitable for
use as a set of token keys with the RSA JWT backend.  The bits
parameter allows you to determine how many bits your keys should be,
values as low as 1024 are technically possible, but 2048 really should
be your minimum for security.
`
)

func init() {
	keygenCmd.AddCommand(keygenRSACmd)
}

func keygenRSACmdRun(cmd *cobra.Command, args []string) {
	bits, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "bits must be a number")
		os.Exit(1)
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating private key: %v", err)
		os.Exit(1)
	}

	pridata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)
	if err := os.WriteFile("rsa-private.tokenkey", pridata, 0400); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing private key: %v", err)
		os.Exit(1)
	}

	pubASN1, _ := x509.MarshalPKIXPublicKey(privateKey.Public())

	pubdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubASN1,
		},
	)
	if err := os.WriteFile("rsa-public.tokenkey", pubdata, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing public key: %v", err)
		os.Exit(1)
	}
	fmt.Println("Keys generated")
}
