package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/NetAuth/NetAuth/internal/token"

	"github.com/dgrijalva/jwt-go"
)

var (
	privateKeyFile = flag.String("jwt_rsa_privatekey", "netauth.key", "Path to private key")
	publicKeyFile  = flag.String("jwt_rsa_publickey", "netauth.pem", "Path to public key")
	rsaBits        = flag.Int("jwt_rsa_bits", 2048, "Bit length of generated keys")
	generate       = flag.Bool("jwt_rsa_generate", false, "Generate keys if not available")
)

// An RSAToken is a token that provides both the token.Claims required
// components and the jtw.StandardClaims.
type RSAToken struct {
	token.Claims
	jwt.StandardClaims
}

// The RSATokenService provides RSA tokens and the means to verify
// them.
type RSATokenService struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func init() {
	token.Register("jwt-rsa", NewRSA)
}

// NewRSA returns an RSATokenService initialized and ready for use.
func NewRSA() (token.Service, error) {
	x := RSATokenService{}
	if err := x.GetKeys(); err != nil {
		return nil, err
	}
	return &x, nil
}

// Generate generates a token signed by an RSA key.
func (s *RSATokenService) Generate(claims token.Claims, config token.Config) (string, error) {
	if s.privateKey == nil {
		// Private key is unavailable, signing is not possible
		return "", token.ErrKeyUnavailable
	}

	claims.RenewalsLeft = config.Renewals
	c := RSAToken{
		claims,
		jwt.StandardClaims{
			IssuedAt:  config.IssuedAt.Unix(),
			NotBefore: config.NotBefore.Unix(),
			ExpiresAt: config.NotBefore.Add(config.Lifetime).Unix(),
			Subject:   "NetAuth Standard Token",
			Audience:  "Unrestricted",
			Issuer:    config.Issuer,
			Id:        claims.EntityID,
		},
	}

	tkn := jwt.NewWithClaims(jwt.SigningMethodRS512, c)
	ss, err := tkn.SignedString(s.privateKey)
	if err != nil {
		return "", token.ErrInternalError
	}
	return ss, nil
}

// Validate validates a token signed by an RSA key.
func (s *RSATokenService) Validate(tkn string) (token.Claims, error) {
	if s.publicKey == nil {
		return token.Claims{}, token.ErrKeyUnavailable
	}

	t, err := jwt.ParseWithClaims(tkn, &RSAToken{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return s.publicKey, nil
	})
	if err != nil {
		// This case gets raised if the token wasn't parsable
		// for some reason, or the signing key was wrong, or
		// it was corrupt in some way.
		return token.Claims{}, token.ErrInternalError
	}

	claims, ok := t.Claims.(*RSAToken)
	if !ok || !t.Valid {
		// This case is raised when the token was parseable,
		// but wasn't validated due to key errors, expiration
		// or other similar problems.
		return token.Claims{}, token.ErrTokenInvalid
	}
	return claims.Claims, nil
}

// GetKeys obtains the keys for an RSATokenService.  If the keys are
// not available and it is not disabled, then a keypair will be
// generated.
func (s *RSATokenService) GetKeys() error {
	log.Printf("Loading public key from %s", *publicKeyFile)
	f, err := ioutil.ReadFile(*publicKeyFile)
	if os.IsNotExist(err) {
		log.Printf("Blob at %s contains no key!", *publicKeyFile)

		if !*generate {
			log.Println("Generating keys is disabled!")
			return token.ErrKeyGenerationDisabled
		}

		log.Println("Generating keys")

		// No keys, we need to create them
		s.privateKey, err = rsa.GenerateKey(rand.Reader, *rsaBits)
		if err != nil {
			log.Println(err)
			return token.ErrInternalError
		}
		s.publicKey = &s.privateKey.PublicKey

		// Marshal the private key
		pridata := pem.EncodeToMemory(
			&pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: x509.MarshalPKCS1PrivateKey(s.privateKey),
			},
		)
		if err := ioutil.WriteFile(*privateKeyFile, pridata, 0400); err != nil {
			log.Println(err)
			return token.ErrInternalError
		}

		// Marshal the public key
		pubASN1, err := x509.MarshalPKIXPublicKey(s.publicKey)
		if err != nil {
			log.Println(err)
			return token.ErrInternalError
		}
		pubdata := pem.EncodeToMemory(
			&pem.Block{
				Type:  "RSA PUBLIC KEY",
				Bytes: pubASN1,
			},
		)
		if err := ioutil.WriteFile(*publicKeyFile, pubdata, 0644); err != nil {
			log.Println(err)
			return token.ErrInternalError
		}
		// At this point the key is saved to disk and
		// initialized
		return nil
	}
	if err != nil {
		log.Println("No key available and generate disabled!")
		return token.ErrKeyUnavailable
	}

	block, _ := pem.Decode([]byte(f))
	if block == nil {
		log.Println("Error decoding PEM block")
		return token.ErrKeyUnavailable
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return token.ErrKeyUnavailable
	}

	p, ok := pub.(*rsa.PublicKey)
	if !ok {
		log.Printf("%s does not contain a public key", *publicKeyFile)
		return token.ErrKeyUnavailable
	}
	s.publicKey = p

	// Now we'll try and load the private key, this doesn't error
	// out, because you can still do meaningful work with the
	// public key.  The generate function will return errors
	// though if the private key fails to load.
	log.Printf("Loading private key from %s", *privateKeyFile)
	pristr, err := ioutil.ReadFile(*privateKeyFile)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("File load error: %s", err)
	}
	if os.IsNotExist(err) {
		// No private key, so we bail out early.  This doesn't
		// return an error because the general case is
		// verifying an existing token, which only needs the
		// public key.  In this case unavailability of the
		// private key will trigger an error on signing.
		log.Println("Token: No private key available, signing will be unavailable")
		return nil
	}

	block, _ = pem.Decode([]byte(pristr))
	if block == nil {
		log.Println("Error decoding PEM block (private key)")
	}
	s.privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// We don't want to error out here since this isn't
		// needed if all you want to do is verify a signature.
		s.privateKey = nil
	}

	// Keys loaded and ready to sign with
	return nil
}
