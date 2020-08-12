package jwt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"

	"github.com/netauth/netauth/internal/token"
)

var (
	config = token.Config{
		Lifetime: time.Minute * 5,
		Issuer:   "NetAuth Test",
	}
)

func init() {
	// We're throwing the error away here since this is parsing
	// the reference format, if that doesn't work its very
	// unlikely anything else does.
	t, _ := time.Parse(time.ANSIC, time.ANSIC)

	config.IssuedAt = t
	config.NotBefore = t

	viper.Set("token.jwt.bits", 1024)
}

func mkTmpTestDir(t *testing.T) string {
	dir, err := ioutil.TempDir("/tmp", "tkntest")
	if err != nil {
		t.Error(err)
	}
	viper.Set("core.home", dir)
	return dir
}

func cleanTmpTestDir(dir string, t *testing.T) {
	// Remove the tmpdir, don't want to clutter the filesystem
	if err := os.RemoveAll(dir); err != nil {
		t.Log(err)
	}
}

func genFixedKey(dir string, t *testing.T) {
	// Chosen by fair dice roll.
	r := rand.New(rand.NewSource(4))

	if err := os.MkdirAll(filepath.Join(dir, "keys"), 0755); err != nil {
		t.Fatal(err)
	}

	// No keys, we need to create them
	privateKey, err := rsa.GenerateKey(r, viper.GetInt("token.jwt.bits"))
	if err != nil {
		t.Fatal("Keys unavailable")
	}
	publicKey := &privateKey.PublicKey

	// Marshal the private key
	pridata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)
	if err := ioutil.WriteFile(filepath.Join(dir, "keys", "token.key"), pridata, 0400); err != nil {
		t.Fatal("Keys unavailable")
	}

	// Marshal the public key
	pubASN1, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		t.Fatal("Keys unavailable")
	}
	pubdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubASN1,
		},
	)
	if err := ioutil.WriteFile(filepath.Join(dir, "keys", "token.pem"), pubdata, 0644); err != nil {
		t.Fatal("Keys unavailable")
	}
}

func TestNewMissingKeys(t *testing.T) {
	testDir := mkTmpTestDir(t)
	defer cleanTmpTestDir(testDir, t)
	viper.Set("token.jwt.generate", false)

	_, err := NewRSA(hclog.NewNullLogger())
	if err != token.ErrKeyGenerationDisabled {
		t.Fatal(err)
	}
}

func TestNewExistingKey(t *testing.T) {
	testDir := mkTmpTestDir(t)
	defer cleanTmpTestDir(testDir, t)

	// This one should generate keys
	viper.Set("token.jwt.generate", true)
	_, err := NewRSA(hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}

	// This one should be loading the existing key
	viper.Set("token.jwt.generate", false)
	_, err = NewRSA(hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateNoKey(t *testing.T) {
	testDir := mkTmpTestDir(t)
	defer cleanTmpTestDir(testDir, t)
	viper.Set("token.jwt.generate", true)

	x, err := NewRSA(hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}

	rx, ok := x.(*RSATokenService)
	if !ok {
		t.Fatal("Bad type assertion")
	}
	rx.privateKey = nil

	if _, err := rx.Generate(token.Claims{}, config); err != token.ErrKeyUnavailable {
		t.Error(err)
	}
}

func TestValidateToken(t *testing.T) {
	testDir := mkTmpTestDir(t)
	defer cleanTmpTestDir(testDir, t)

	// generate a fixed value key
	genFixedKey(testDir, t)

	// Create the token service which will use the key generated
	// earlier
	x, err := NewRSA(hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}

	c := token.Claims{
		EntityID: "foo",
	}

	cfg := token.Config{
		Lifetime:  time.Minute * 5,
		IssuedAt:  time.Now(),
		NotBefore: time.Now(),
		Issuer:    "NetAuth Test",
	}

	tkn, err := x.Generate(c, cfg)
	if err != nil {
		t.Fatal(err)
	}

	os.Remove(filepath.Join(testDir, "keys", "token.key"))
	x, err = NewRSA(hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}

	claims, err := x.Validate(tkn)
	if err != nil {
		t.Fatal(err)
	}

	// structs containing []string can't be compared directly, so
	// we compare the value that was set earlier
	if claims.EntityID != c.EntityID {
		t.Error("Claims are not the same!")
	}
}

func TestValidateNoKey(t *testing.T) {
	testDir := mkTmpTestDir(t)
	defer cleanTmpTestDir(testDir, t)

	// generate a fixed value key
	genFixedKey(testDir, t)

	// Create the token service which will use the key generated
	// earlier
	x, err := NewRSA(hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}

	rx, ok := x.(*RSATokenService)
	if !ok {
		t.Fatal("Type Error")
	}

	rx.publicKey = nil

	if _, err := x.Validate(""); err != token.ErrKeyUnavailable {
		t.Error(err)
	}
}

func TestValidateCorruptToken(t *testing.T) {
	testDir := mkTmpTestDir(t)
	defer cleanTmpTestDir(testDir, t)

	// generate a fixed value key
	genFixedKey(testDir, t)

	// Create the token service which will use the key generated
	// earlier
	x, err := NewRSA(hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}

	if _, err := x.Validate(""); err != token.ErrInternalError {
		t.Error(err)
	}
}

func TestValidateWrongSigningMethod(t *testing.T) {
	testDir := mkTmpTestDir(t)
	defer cleanTmpTestDir(testDir, t)

	// generate a fixed value key
	genFixedKey(testDir, t)

	// Create the token service which will use the key generated
	// earlier
	x, err := NewRSA(hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}

	badToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	if _, err := x.Validate(badToken); err != token.ErrTokenInvalid {
		t.Logf("%T", err)
		t.Error(err)
	}
}

func TestValidateExpiredToken(t *testing.T) {
	testDir := mkTmpTestDir(t)
	defer cleanTmpTestDir(testDir, t)

	// generate a fixed value key
	genFixedKey(testDir, t)

	// Create the token service which will use the key generated
	// earlier
	x, err := NewRSA(hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}

	c := token.Claims{
		EntityID: "foo",
	}

	cfg := token.Config{
		Lifetime:  0,
		IssuedAt:  time.Now().Add(-1 * time.Minute),
		NotBefore: time.Now().Add(-1 * time.Minute),
		Issuer:    "NetAuth Test",
	}

	tkn, err := x.Generate(c, cfg)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := x.Validate(tkn); err != token.ErrTokenInvalid {
		t.Logf("%T", err)
		t.Log(err)
	}
}

func TestGetKeysNoGenerate(t *testing.T) {
	testDir := mkTmpTestDir(t)
	defer cleanTmpTestDir(testDir, t)
	viper.Set("token.jwt.generate", true)

	if err := os.MkdirAll(filepath.Join(testDir, "keys", "token.key"), 0755); err != nil {
		t.Fatal(err)
	}

	_, err := NewRSA(hclog.NewNullLogger())
	if err != token.ErrInternalError {
		t.Error(err)
	}
}

func TestGetKeysBadPublicKeyFile(t *testing.T) {
	testDir := mkTmpTestDir(t)
	defer cleanTmpTestDir(testDir, t)
	viper.Set("token.jwt.generate", false)

	if err := os.MkdirAll(filepath.Join(testDir, "keys", "token.pem"), 0755); err != nil {
		t.Fatal(err)
	}

	_, err := NewRSA(hclog.NewNullLogger())
	if err != token.ErrKeyUnavailable {
		t.Error(err)
	}

}

func TestGetKeysBadPublicKeyMode(t *testing.T) {
	testDir := mkTmpTestDir(t)
	defer cleanTmpTestDir(testDir, t)

	genFixedKey(testDir, t)

	if err := os.Chmod(filepath.Join(testDir, "keys", "token.pem"), 0400); err != nil {
		t.Fatal(err)
	}

	_, err := NewRSA(hclog.NewNullLogger())
	if err != token.ErrKeyUnavailable {
		t.Error(err)
	}
}

func TestGetKeysBadBlockDecode(t *testing.T) {
	testDir := mkTmpTestDir(t)
	defer cleanTmpTestDir(testDir, t)
	viper.Set("token.jwt.generate", true)

	if err := os.MkdirAll(filepath.Join(testDir, "keys"), 0755); err != nil {
		t.Fatal(err)
	}

	if _, err := os.OpenFile(filepath.Join(testDir, "keys", "token.pem"), os.O_RDONLY|os.O_CREATE, 0666); err != nil {
		t.Fatal(err)
	}

	_, err := NewRSA(hclog.NewNullLogger())
	if err != token.ErrKeyUnavailable {
		t.Error(err)
	}
}

func TestGetKeysPublicKeyWrongType(t *testing.T) {
	testDir := mkTmpTestDir(t)
	defer cleanTmpTestDir(testDir, t)
	viper.Set("token.jwt.generate", false)

	// Chosen by fair dice roll.
	r := rand.New(rand.NewSource(4))
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), r)
	if err != nil {
		t.Fatal(err)
	}

	// Marshal the public key
	pubASN1, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Log(err)
		t.Fatal("Couldn't marshal key")
	}
	pubdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: pubASN1,
		},
	)

	if err := os.MkdirAll(filepath.Join(testDir, "keys"), 0755); err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile(filepath.Join(testDir, "keys", "token.pem"), pubdata, 0644); err != nil {
		t.Fatal("Keys unavailable")
	}

	_, err = NewRSA(hclog.NewNullLogger())
	if err != token.ErrKeyUnavailable {
		t.Error(err)
	}
}

func TestGetKeysPublicKeyIsPrivate(t *testing.T) {
	testDir := mkTmpTestDir(t)
	defer cleanTmpTestDir(testDir, t)
	viper.Set("token.jwt.generate", false)

	// Generate the keys flipped, so that the private key winds up
	// in the wrong file, then flip the keys back.
	genFixedKey(testDir, t)

	b := filepath.Join(testDir, "keys")
	if err := os.Remove(filepath.Join(b, "token.pem")); err != nil {
		t.Fatal(err)
	}
	if err := os.Link(filepath.Join(b, "token.key"), filepath.Join(b, "token.pem")); err != nil {
		t.Fatal(err)
	}

	if err := os.Chmod(filepath.Join(b, "token.key"), 0644); err != nil {
		t.Fatal("Couldn't set mode on key")
	}

	_, err := NewRSA(hclog.NewNullLogger())
	if err != token.ErrKeyUnavailable {
		t.Error(err)
	}
}

func TestGetKeysNoPrivateKey(t *testing.T) {
	testDir := mkTmpTestDir(t)
	defer cleanTmpTestDir(testDir, t)
	viper.Set("token.jwt.generate", false)

	genFixedKey(testDir, t)
	if err := os.Remove(filepath.Join(testDir, "keys", "token.key")); err != nil {
		t.Fatal(err)
	}

	_, err := NewRSA(hclog.NewNullLogger())
	if err != nil {
		t.Error(err)
	}
}

func TestGetKeysUnreadablePrivateKey(t *testing.T) {
	testDir := mkTmpTestDir(t)
	defer cleanTmpTestDir(testDir, t)
	viper.Set("token.jwt.generate", false)

	genFixedKey(testDir, t)

	if err := os.Chmod(filepath.Join(testDir, "keys", "token.key"), 0000); err != nil {
		t.Fatal(err)
	}

	_, err := NewRSA(hclog.NewNullLogger())
	if err != nil {
		t.Error(err)
	}
}

func TestGetKeysPrivateKeyIsPublic(t *testing.T) {
	testDir := mkTmpTestDir(t)
	defer cleanTmpTestDir(testDir, t)
	viper.Set("token.jwt.generate", false)

	// Write out a public key, then change where the private key
	// points to cause a load error.  This must not cause an
	// error.
	genFixedKey(testDir, t)

	b := filepath.Join(testDir, "keys")
	if err := os.Remove(filepath.Join(b, "token.key")); err != nil {
		t.Fatal(err)
	}
	if err := os.Link(filepath.Join(b, "token.pem"), filepath.Join(b, "token.key")); err != nil {
		t.Fatal(err)
	}

	_, err := NewRSA(hclog.NewNullLogger())
	if err != nil {
		t.Error(err)
	}
}

func TestGenerateKeysSuccess(t *testing.T) {
	testDir := mkTmpTestDir(t)
	defer cleanTmpTestDir(testDir, t)
	viper.Set("token.jwt.generate", true)

	x, err := NewRSA(hclog.NewNullLogger())
	if err != nil {
		t.Error(err)
	}

	rx, ok := x.(*RSATokenService)
	if !ok {
		t.Fatal("Type Error")
	}

	if err := os.RemoveAll(filepath.Join(testDir, "keys")); err != nil {
		t.Fatal(err)
	}

	if err := rx.generateKeys(256); err != nil {
		t.Error(err)
	}
}

func TestGenerateKeysWrongBitNumber(t *testing.T) {
	testDir := mkTmpTestDir(t)
	defer cleanTmpTestDir(testDir, t)
	viper.Set("token.jwt.generate", true)

	x, err := NewRSA(hclog.NewNullLogger())
	if err != nil {
		t.Error(err)
	}

	rx, ok := x.(*RSATokenService)
	if !ok {
		t.Fatal("Type Error")
	}

	if err := rx.generateKeys(0); err != token.ErrInternalError {
		t.Error(err)
	}
}

func TestGenerateKeysBadPrivateKeyFile(t *testing.T) {
	testDir := mkTmpTestDir(t)
	defer cleanTmpTestDir(testDir, t)
	viper.Set("token.jwt.generate", true)

	x, err := NewRSA(hclog.NewNullLogger())
	if err != nil {
		t.Error(err)
	}

	rx, ok := x.(*RSATokenService)
	if !ok {
		t.Fatal("Type Error")
	}

	pkf := filepath.Join(testDir, "keys", "token.key")

	if err := os.Remove(pkf); err != nil {
		t.Fatal(err)
	}

	if err := os.Mkdir(pkf, 0755); err != nil {
		t.Fatal(err)
	}

	if err := rx.generateKeys(256); err != token.ErrInternalError {
		t.Error(err)
	}
}

func TestGenerateKeysBadPublicKeyFile(t *testing.T) {
	testDir := mkTmpTestDir(t)
	defer cleanTmpTestDir(testDir, t)
	viper.Set("token.jwt.generate", true)

	x, err := NewRSA(hclog.NewNullLogger())
	if err != nil {
		t.Error(err)
	}

	rx, ok := x.(*RSATokenService)
	if !ok {
		t.Fatal("Type Error")
	}

	// Remove the private key file since the system tries to write
	// it first.
	if err := os.Remove(filepath.Join(testDir, "keys", "token.key")); err != nil {
		t.Fatal(err)
	}

	pkf := filepath.Join(testDir, "keys", "badPublicKey")
	if err := os.Mkdir(pkf, 0755); err != nil {
		t.Fatal(err)
	}

	rx.publicKeyFile = pkf

	if err := rx.generateKeys(256); err != token.ErrInternalError {
		t.Error(err)
	}
}

func TestHealthCheck(t *testing.T) {
	testDir := mkTmpTestDir(t)
	defer cleanTmpTestDir(testDir, t)
	viper.Set("token.jwt.generate", true)

	x, err := NewRSA(hclog.NewNullLogger())
	if err != nil {
		t.Error(err)
	}

	rx, ok := x.(*RSATokenService)
	if !ok {
		t.Fatal("Type Error")
	}

	if status := rx.healthCheck(); !status.OK {
		t.Error(status)
	}
}

func TestHealthCheckNoPublicKey(t *testing.T) {
	testDir := mkTmpTestDir(t)
	defer cleanTmpTestDir(testDir, t)
	viper.Set("token.jwt.generate", true)

	x, err := NewRSA(hclog.NewNullLogger())
	if err != nil {
		t.Error(err)
	}

	rx, ok := x.(*RSATokenService)
	if !ok {
		t.Fatal("Type Error")
	}

	rx.publicKey = nil

	if status := rx.healthCheck(); status.OK {
		t.Error(status)
	}
}

func TestHealthCheckNoPrivateKey(t *testing.T) {
	testDir := mkTmpTestDir(t)
	defer cleanTmpTestDir(testDir, t)
	viper.Set("token.jwt.generate", true)

	x, err := NewRSA(hclog.NewNullLogger())
	if err != nil {
		t.Error(err)
	}

	rx, ok := x.(*RSATokenService)
	if !ok {
		t.Fatal("Type Error")
	}

	rx.privateKey = nil

	if status := rx.healthCheck(); status.OK {
		t.Error(status)
	}
}

func TestHealthCheckBadPrivateKeyPermissions(t *testing.T) {
	testDir := mkTmpTestDir(t)
	defer cleanTmpTestDir(testDir, t)
	viper.Set("token.jwt.generate", true)

	x, err := NewRSA(hclog.NewNullLogger())
	if err != nil {
		t.Error(err)
	}

	if err := os.Chmod(filepath.Join(testDir, "keys", "token.key"), 0644); err != nil {
		t.Fatal(err)
	}

	rx, ok := x.(*RSATokenService)
	if !ok {
		t.Fatal("Type Error")
	}

	if status := rx.healthCheck(); status.OK {
		t.Error(status)
	}
}

func TestHealthCheckBadPublicKeyPermissions(t *testing.T) {
	testDir := mkTmpTestDir(t)
	defer cleanTmpTestDir(testDir, t)
	viper.Set("token.jwt.generate", true)

	x, err := NewRSA(hclog.NewNullLogger())
	if err != nil {
		t.Error(err)
	}

	if err := os.Chmod(filepath.Join(testDir, "keys", "token.pem"), 0600); err != nil {
		t.Fatal(err)
	}

	rx, ok := x.(*RSATokenService)
	if !ok {
		t.Fatal("Type Error")
	}

	if status := rx.healthCheck(); status.OK {
		t.Error(status)
	}
}

func TestCheckKeyModeOKBadStat(t *testing.T) {
	testDir := mkTmpTestDir(t)
	defer cleanTmpTestDir(testDir, t)
	viper.Set("token.jwt.generate", true)

	x, err := NewRSA(hclog.NewNullLogger())
	if err != nil {
		t.Error(err)
	}

	rx, ok := x.(*RSATokenService)
	if !ok {
		t.Fatal("Type Error")
	}

	if rx.checkKeyModeOK("", filepath.Join(testDir, "does-not-exist")) {
		t.Error("Stat succeeded on a non-existent path")
	}
}
