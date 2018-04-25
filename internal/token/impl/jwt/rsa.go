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
	rsa_bits       = flag.Int("jwt_rsa_bits", 2048, "Bit length of generated keys")
	generate       = flag.Bool("jwt_rsa_generate", false, "Generate keys if not available")
)

type RSAToken struct {
	token.Claims
	jwt.StandardClaims
}

type RSATokenService struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func init() {
	token.RegisterService("jwt-rsa", NewRSA)
}

func NewRSA() (token.TokenService, error) {
	x := RSATokenService{}
	if err := x.GetKeys(); err != nil {
		return nil, err
	}
	return &x, nil
}

func (s *RSATokenService) Generate(claims token.Claims, config token.TokenConfig) (string, error) {
	if s.privateKey == nil {
		// Private key is unavailable, signing is not possible
		return "", token.KEY_UNAVAILABLE
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
		return "", err
	}
	return ss, nil
}

func (s *RSATokenService) Validate(tkn string) (token.Claims, error) {
	if s.publicKey == nil {
		return token.Claims{}, token.KEY_UNAVAILABLE
	}

	t, err := jwt.ParseWithClaims(tkn, &RSAToken{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return s.publicKey, nil
	})
	if err != nil {
		return token.Claims{}, err
	}

	if claims, ok := t.Claims.(*RSAToken); ok && t.Valid {
		return claims.Claims, err
	} else {
		return token.Claims{}, err
	}
}

func (s *RSATokenService) GetKeys() error {
	log.Printf("Loading public key from %s", *publicKeyFile)
	f, err := ioutil.ReadFile(*publicKeyFile)
	if os.IsNotExist(err) {
		log.Printf("Blob at %s contains no key!", *publicKeyFile)

		if !*generate {
			log.Println("Generating keys is disabled!")
			return token.NO_GENERATE_KEYS
		}

		log.Println("Generating keys")

		// No keys, we need to create them
		s.privateKey, err = rsa.GenerateKey(rand.Reader, *rsa_bits)
		if err != nil {
			log.Println(err)
			return token.KEY_UNAVAILABLE
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
			return token.KEY_UNAVAILABLE
		}

		// Marshal the public key
		pubASN1, err := x509.MarshalPKIXPublicKey(s.publicKey)
		if err != nil {
			log.Println(err)
			return token.KEY_UNAVAILABLE
		}
		pubdata := pem.EncodeToMemory(
			&pem.Block{
				Type:  "RSA PUBLIC KEY",
				Bytes: pubASN1,
			},
		)
		if err := ioutil.WriteFile(*publicKeyFile, pubdata, 0644); err != nil {
			log.Println(err)
			return token.KEY_UNAVAILABLE
		}
		// At this point the key is saved to disk and
		// initialized
		return nil
	}
	if err != nil {
		log.Println("No key available and generate disabled!")
		return token.KEY_UNAVAILABLE
	}

	block, _ := pem.Decode([]byte(f))
	if block == nil {
		log.Println("Error decoding PEM block")
		return token.KEY_UNAVAILABLE
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return token.KEY_UNAVAILABLE
	}

	p, ok := pub.(*rsa.PublicKey)
	if !ok {
		log.Printf("%s does not contain a public key", *publicKeyFile)
		return token.KEY_UNAVAILABLE
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
		// No private key, so we bail out early
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
