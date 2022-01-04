package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/hashicorp/go-hclog"

	"github.com/netauth/netauth/pkg/token"
	"github.com/netauth/netauth/pkg/token/keyprovider"

	"github.com/golang-jwt/jwt/v4"
)

// An RSAToken is a token that provides both the token.Claims required
// components and the jtw.StandardClaims.
type RSAToken struct {
	token.Claims
	jwt.RegisteredClaims
}

// The RSATokenService provides RSA tokens and the means to verify
// them.
type RSATokenService struct {
	log hclog.Logger

	kp keyprovider.KeyProvider
}

func init() {
	token.Register("jwt-rsa", NewRSA)
}

// NewRSA returns an RSATokenService initialized and ready for use.
func NewRSA(l hclog.Logger, kp keyprovider.KeyProvider) (token.Service, error) {
	x := RSATokenService{}
	x.log = l.Named("jwt-rsa")
	x.kp = kp

	return &x, nil
}

// Generate generates a token signed by an RSA key.
func (s *RSATokenService) Generate(claims token.Claims, config token.Config) (string, error) {
	key, err := s.privkey()
	if err != nil {
		return "", err
	}

	c := RSAToken{
		claims,
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(config.IssuedAt),
			NotBefore: jwt.NewNumericDate(config.NotBefore),
			ExpiresAt: jwt.NewNumericDate(config.NotBefore.Add(config.Lifetime)),
			Subject:   "netauth/" + claims.EntityID,
			Audience:  []string{"netauth-internal"},
			Issuer:    config.Issuer,
			ID:        claims.EntityID,
		},
	}

	tkn := jwt.NewWithClaims(jwt.SigningMethodRS512, c)

	// We discard this error as there is no meaningful error that
	// can be returned from here.  Basically the FPU would need to
	// fail for this to have a problem...
	ss, _ := tkn.SignedString(key)
	return ss, nil
}

// Validate validates a token signed by an RSA key.
func (s *RSATokenService) Validate(tkn string) (token.Claims, error) {
	k, err := s.pubkey()
	if err != nil {
		return token.Claims{}, err
	}

	t, err := jwt.ParseWithClaims(tkn, &RSAToken{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			s.log.Error("Token was signed with invalid algorithm", "expected", t.Method, "actual", t.Header["alg"])
			return nil, token.ErrTokenInvalid
		}
		return k, nil
	})
	if err != nil {
		// This case gets raised if the token wasn't parsable
		// for some reason, or the signing key was wrong, or
		// it was corrupt in some way.
		if t != nil && !t.Valid {
			return token.Claims{}, token.ErrTokenInvalid
		}
		return token.Claims{}, token.ErrInternalError
	}

	// We do a blind type change here to pull out the embedded
	// RSAToken which includes a token.Claims.  We can be sure
	// this is an RSAToken because if it wasn't, the
	// ParseWithClaims call would have exploded just above.
	claims, _ := t.Claims.(*RSAToken)
	return claims.Claims, nil
}

func (s *RSATokenService) pubkey() (*rsa.PublicKey, error) {
	b, err := s.kp.Provide("rsa", "public")
	switch err {
	case nil:
		break
	case keyprovider.ErrNoSuchKey:
		return nil, token.ErrKeyUnavailable
	default:
		return nil, err
	}
	block, _ := pem.Decode(b)
	if block == nil {
		s.log.Error("Error decoding PEM block")
		return nil, token.ErrKeyUnavailable
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		s.log.Error("Error parsing key:", "error", err)
		return nil, token.ErrKeyUnavailable
	}

	p, ok := pub.(*rsa.PublicKey)
	if !ok {
		s.log.Error("Parsed key is not an RSA Public Key")
		return nil, token.ErrKeyUnavailable
	}
	return p, nil
}

func (s *RSATokenService) privkey() (*rsa.PrivateKey, error) {
	b, err := s.kp.Provide("rsa", "private")
	switch err {
	case nil:
		break
	case keyprovider.ErrNoSuchKey:
		return nil, token.ErrKeyUnavailable
	default:
		return nil, err
	}
	block, _ := pem.Decode(b)
	if block == nil {
		return nil, token.ErrKeyUnavailable
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}
