package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	certgenCmd = &cobra.Command{
		Use:   "certgen <cn> [san1]..[sanN]",
		Short: "Generate a self-signed ecdsa certificate for the provided names",
		Long:  certgenCmdLongDocs,
		Run:   certgenCmdRun,
		Args:  cobra.MinimumNArgs(1),
	}

	certgenCmdLongDocs = `The certgen command will create a certificate and key that are
suitable for configuring your netauth environment for a demo or
testing.  The generated certificate is not suitable for long term or
production use.`

	certgenCmdValidityDuration   time.Duration
	certgenCmdIssuerOrganization string
)

func init() {
	rootCmd.AddCommand(certgenCmd)

	certgenCmd.Flags().StringVarP(&certgenCmdIssuerOrganization, "organization", "o", "nsutil", "Certificate Issuer Organization")
	certgenCmd.Flags().DurationVarP(&certgenCmdValidityDuration, "validity", "d", time.Hour*24*180, "Validity of the generated certificate")
}

func certgenCmdRun(cmd *cobra.Command, args []string) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating keys: %s\n", err)
		os.Exit(1)
	}
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{certgenCmdIssuerOrganization},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(certgenCmdValidityDuration),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	for _, h := range args {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	b, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
		os.Exit(2)
	}
	pb := &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create certificate: %s\n", err)
		os.Exit(2)
	}
	out := &bytes.Buffer{}
	pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	if err := os.WriteFile("tls.pem", out.Bytes(), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing certificate: %s\n", err)
		os.Exit(3)
	}

	if err := os.WriteFile("tls.key", pem.EncodeToMemory(pb), 0400); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing key: %s\n", err)
		os.Exit(3)
	}
}
