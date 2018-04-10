package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/NetAuth/NetAuth/internal/token"

	"github.com/dgrijalva/jwt-go"
)

var (
	key_blob = flag.String("jwt_rsa_key", "key.dat", "Path to key file")
	rsa_bits = flag.Int("jwt_rsa_bits", 2048, "Bit length of generated keys")
	generate = flag.Bool("jwt_rsa_generate", true, "Generate keys if not available")
)

type RSAToken struct {
	token.Claims
	jwt.StandardClaims
}

type RSATokenService struct {
	key *rsa.PrivateKey
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
	ss, err := tkn.SignedString(s.key)
	if err != nil {
		return "", err
	}
	return ss, nil
}

func (s *RSATokenService) Validate(tkn string) (token.Claims, error) {
	t, err := jwt.ParseWithClaims(tkn, &RSAToken{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return s.key.Public(), nil
	})

	if claims, ok := t.Claims.(*RSAToken); ok && t.Valid {
		return claims.Claims, err
	} else {
		return token.Claims{}, err
	}
}

func (s *RSATokenService) GetKeys() error {
	log.Printf("Loading key blob from %s", *key_blob)
	f, err := ioutil.ReadFile(*key_blob)
	if os.IsNotExist(err) {
		log.Printf("Blob at %s contains no key!", *key_blob)

		if !*generate {
			log.Println("Generating keys is disabled!")
			return token.NO_GENERATE_KEYS
		}

		log.Println("Generating keys")

		// No keys, we need to create them
		s.key, err = rsa.GenerateKey(rand.Reader, *rsa_bits)
		if err != nil {
			log.Println(err)
			return token.KEY_UNAVAILABLE
		}
		d, err := json.Marshal(s.key)
		if err != nil {
			log.Println(err)
			return token.KEY_UNAVAILABLE
		}
		if err := ioutil.WriteFile(*key_blob, d, 0400); err != nil {
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

	if err := json.Unmarshal(f, &s.key); err != nil {
		return err
	}

	// Keys loaded and ready to sign with
	return nil
}
