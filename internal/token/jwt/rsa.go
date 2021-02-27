package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"

	"github.com/netauth/netauth/internal/health"
	"github.com/netauth/netauth/internal/token"

	"github.com/dgrijalva/jwt-go"
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

	publicKeyFile  string
	privateKeyFile string

	log hclog.Logger
}

func init() {
	token.Register("jwt-rsa", NewRSA)
}

// NewRSA returns an RSATokenService initialized and ready for use.
func NewRSA(l hclog.Logger) (token.Service, error) {
	x := RSATokenService{}
	x.log = l.Named("jwt-rsa")

	if err := x.GetKeys(); err != nil {
		return nil, err
	}

	health.RegisterCheck("JWT-RSA", x.healthCheck)

	return &x, nil
}

// Generate generates a token signed by an RSA key.
func (s *RSATokenService) Generate(claims token.Claims, config token.Config) (string, error) {
	if s.privateKey == nil {
		// Private key is unavailable, signing is not possible
		return "", token.ErrKeyUnavailable
	}

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

	// We discard this error as there is no meaningful error that
	// can be returned from here.  Basically the FPU would need to
	// fail for this to have a problem...
	ss, _ := tkn.SignedString(s.privateKey)
	return ss, nil
}

// Validate validates a token signed by an RSA key.
func (s *RSATokenService) Validate(tkn string) (token.Claims, error) {
	if s.publicKey == nil {
		return token.Claims{}, token.ErrKeyUnavailable
	}

	t, err := jwt.ParseWithClaims(tkn, &RSAToken{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			s.log.Error("Token was signed with invalid algorithm", "expected", t.Method, "actual", t.Header["alg"])
			return nil, token.ErrTokenInvalid
		}
		return s.publicKey, nil
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

// GetKeys obtains the keys for an RSATokenService.  If the keys are
// not available and it is not disabled, then a keypair will be
// generated.
func (s *RSATokenService) GetKeys() error {
	s.publicKeyFile = filepath.Join(viper.GetString("core.conf"), "keys", "token.pem")
	s.privateKeyFile = filepath.Join(viper.GetString("core.conf"), "keys", "token.key")

	s.log.Debug("Loading public key", "file", s.publicKeyFile)
	f, err := ioutil.ReadFile(s.publicKeyFile)
	if os.IsNotExist(err) {
		s.log.Error("File contains no key!", "file", s.publicKeyFile)

		if !viper.GetBool("token.jwt.generate") {
			s.log.Warn("Key generation is disabled")
			return token.ErrKeyGenerationDisabled
		}
		s.log.Info("Generating keys")

		// Request the keys be generated
		if err := s.generateKeys(viper.GetInt("token.jwt.bits")); err != nil {
			s.log.Error("Error generating keys", "error", err)
			return err
		}

		// Keys are generated, return out
		return nil
	}
	if err != nil && !os.IsNotExist(err) {
		s.log.Warn("Keys are not available")
		return token.ErrKeyUnavailable
	}

	if !s.checkKeyModeOK("-rw-r--r--", s.publicKeyFile) {
		s.log.Warn("Public key has incorrect mode bits")
		return token.ErrKeyUnavailable
	}

	block, _ := pem.Decode([]byte(f))
	if block == nil {
		s.log.Error("Error decoding PEM block")
		return token.ErrKeyUnavailable
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		s.log.Error("Error parsing key:", "key", s.publicKeyFile, "error", err)
		return token.ErrKeyUnavailable
	}

	p, ok := pub.(*rsa.PublicKey)
	if !ok {
		s.log.Error("File does not contain an RSA public key", "file", s.publicKeyFile)
		return token.ErrKeyUnavailable
	}
	s.publicKey = p

	// Now we'll try and load the private key, this doesn't error
	// out, because you can still do meaningful work with the
	// public key.  The generate function will return errors
	// though if the private key fails to load.
	s.log.Debug("Loading private key from file", "file", s.privateKeyFile)
	pristr, err := ioutil.ReadFile(s.privateKeyFile)
	if err != nil && !os.IsNotExist(err) {
		s.log.Error("Private key load error", "file", s.privateKeyFile, "error", err)
	}
	if os.IsNotExist(err) {
		// No private key, so we bail out early.  This doesn't
		// return an error because the general case is
		// verifying an existing token, which only needs the
		// public key.  In this case unavailability of the
		// private key will trigger an error on signing.
		s.log.Warn("No private key is loaded")
		s.log.Warn("Signing will be unavailable")
		return nil
	}

	if !s.checkKeyModeOK("-r--------", s.privateKeyFile) {
		s.log.Warn("Private key has incorrect mode bits", "file", s.privateKeyFile)
	}

	block, _ = pem.Decode([]byte(pristr))
	if block == nil {
		// We don't want to error out here since this isn't
		// needed if all you want to do is verify a signature.
		s.privateKey = nil
		s.log.Warn("Error decoding private key", "file", s.privateKeyFile)
		return nil
	}
	s.privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// We don't want to error out here since this isn't
		// needed if all you want to do is verify a signature.
		s.privateKey = nil
		s.log.Warn("Error unmarshaling private key", "file", s.privateKeyFile, "error", err)
	}

	// Keys loaded and ready to sign with
	return nil
}

func (s *RSATokenService) generateKeys(bits int) error {
	s.log.Debug("Generating keys")

	// First create the directory for the keys if it doesn't
	// already exist.
	path := filepath.Join(viper.GetString("core.conf"), "keys")
	if err := os.MkdirAll(path, 0755); err != nil {
		s.log.Error("Could not create key directory", "path", path)
		return token.ErrInternalError
	}

	// No keys, we need to create them
	var err error
	s.privateKey, err = rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		s.log.Error("Error generating keys", "error", err)
		return token.ErrInternalError
	}
	s.publicKey = &s.privateKey.PublicKey

	if err := marshalPrivateKey(s.privateKey, s.privateKeyFile); err != nil {
		return err
	}

	if err := marshalPublicKey(s.publicKey, s.publicKeyFile); err != nil {
		return err
	}

	// At this point the key is saved to disk and
	// initialized
	s.log.Debug("Finished generating keys")
	return nil
}

// healthCheck provides a sanity check that keys are loaded and owned
// correctly.
func (s *RSATokenService) healthCheck() health.SubsystemStatus {
	status := health.SubsystemStatus{
		OK:   false,
		Name: "TKN_JWT-RSA",
	}

	if s.privateKey == nil {
		status.Status = "No private key is loaded"
		return status
	}

	if s.publicKey == nil {
		status.Status = "No public key is loaded"
		return status
	}

	if !s.checkKeyModeOK("-rw-r--r--", s.publicKeyFile) {
		status.Status = "Public key has incorrect mode"
		return status
	}

	if !s.checkKeyModeOK("-r--------", s.privateKeyFile) {
		status.Status = "Private key has incorrect mode"
		return status
	}

	status.OK = true
	status.Status = "JWT-RSA TokenService is ready to issue/verify tokens"

	return status
}

func (s *RSATokenService) checkKeyModeOK(mode string, path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		s.log.Error("Error stating key", "error", err)
		return false
	}
	if stat.Mode().Perm().String() != mode {
		s.log.Error("Key permissions are wrong.", "current", stat.Mode().Perm(), "want", mode)
		return false
	}
	return true
}

func marshalPrivateKey(key *rsa.PrivateKey, path string) error {
	pridata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)
	if err := ioutil.WriteFile(path, pridata, 0400); err != nil {
		return token.ErrInternalError
	}

	return nil
}

func marshalPublicKey(key *rsa.PublicKey, path string) error {
	// This error is discarded as there is no case where a
	// meaningful error can be returned from this function that
	// would not already have been caught while marshaling the
	// private key.
	pubASN1, _ := x509.MarshalPKIXPublicKey(key)

	pubdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubASN1,
		},
	)
	if err := ioutil.WriteFile(path, pubdata, 0644); err != nil {
		return token.ErrInternalError
	}
	return nil
}
