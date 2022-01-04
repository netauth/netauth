package jwt

import (
	"errors"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"

	"github.com/netauth/netauth/pkg/token"
	"github.com/netauth/netauth/pkg/token/keyprovider"
	"github.com/netauth/netauth/pkg/token/keyprovider/mock"
)

const (
	pubkey1 = `-----BEGIN RSA PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDLDB82io6KJO2eDdagHnwMxt6m
eA7Fuc2TxeZM6pzb/2W4+wpkmBwwmwfQpH9BK1MWHD2NNS5e7XkDU+c4ja70a6MV
xuztu4YD3kJrHDs1j7BUtlHOM2y1OXBIBG7Cg/BetiTE2Yb5/xS2VgA1wiHrr0M6
3Dt8Rb0D3+5o9ak2yQIDAQAB
-----END RSA PUBLIC KEY-----`

	privkey1 = `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDLDB82io6KJO2eDdagHnwMxt6meA7Fuc2TxeZM6pzb/2W4+wpk
mBwwmwfQpH9BK1MWHD2NNS5e7XkDU+c4ja70a6MVxuztu4YD3kJrHDs1j7BUtlHO
M2y1OXBIBG7Cg/BetiTE2Yb5/xS2VgA1wiHrr0M63Dt8Rb0D3+5o9ak2yQIDAQAB
AoGAQqk8Jh/fJCNzj4xjhjX77AXuWyDXWLrjbzxtm5r63I9AyjZA9z2pI5wCONGI
pdCfeobS/mUTUD8Ol7UYGE0LvsUEoPh81x4QqLb734VKtWRbzEi1PDGX4z/DdD1q
EU1HjrLMw5TgOGne/AMp8pmULC8mhoEI0BszIEqrfjKHxcECQQDzIRAHGWkI//k/
/oVwMcaiF//CidgHQuGAgGCgmz5CESjPslD510jgzFiOhCdNaSkbZ3zv95d6fTy9
EcBwfmhVAkEA1cvcwkeNtJKe01LoFdskOBApforv86uN3FyCh/gVO1dt76OKLNYJ
PBUuluq8USMbufKdO2tt9JGMPi6+uMgYpQJAR0inWV2C5UefvbqTRxzg/z+IFnKx
6xcZ5MI/EnfR3i8HxzWh9k6/qGFhiY+HsnOlwMor4HO4bwpvF4Qv5wu47QJBAK5F
b9yZkOPpRDfD89Sk/eAJJJm2zSNV6tv+OJR232+ws7dMGnyzt3FXXtO74edNc/Nd
1VazGjzqS2QAnIxo5tUCQAw+4zLMuGn6pTQ+eNx93w6xOxjPO/JPZpjrCytZu3Du
oAFH9o0Tnfyf+4C0fpetqJ6PBP4a8tDhW7/88HF03to=
-----END RSA PRIVATE KEY-----`

	pubkey2 = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEpKBU2N2kK2VAaW3Uo/qs69S0A6IT
CbzoF2ZT2ttI6BCiZoLTX2Au9cWFtUUyCEWM+amY9SK3RCkIxCXnBnopYA==
-----END PUBLIC KEY-----`

	pubkey3 = `-----BEGIN RSA PUBLIC KEY-----
aGVsbG8sIHdoeSBhcmUgeW91IGhlcmU/
-----END RSA PUBLIC KEY-----
`
)

func TestGenerateBadKey(t *testing.T) {
	m, _ := mock.New(hclog.NewNullLogger())
	x, err := NewRSA(hclog.NewNullLogger(), m)
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

	m.(*mock.Provider).On("Provide", "rsa", "private").Return([]byte(nil), keyprovider.ErrNoSuchKey)
	if _, err := x.Generate(c, cfg); err != token.ErrKeyUnavailable {
		t.Fatal(err)
	}
}

func TestValidateToken(t *testing.T) {
	m, _ := mock.New(hclog.NewNullLogger())
	x, err := NewRSA(hclog.NewNullLogger(), m)
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

	m.(*mock.Provider).On("Provide", "rsa", "private").Return([]byte(privkey1), nil)
	tkn, err := x.Generate(c, cfg)
	if err != nil {
		t.Fatal(err)
	}

	x, err = NewRSA(hclog.NewNullLogger(), m)
	if err != nil {
		t.Fatal(err)
	}

	m.(*mock.Provider).On("Provide", "rsa", "public").Return([]byte(pubkey1), nil)
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
	m, _ := mock.New(hclog.NewNullLogger())
	x, err := NewRSA(hclog.NewNullLogger(), m)
	if err != nil {
		t.Fatal(err)
	}
	m.(*mock.Provider).On("Provide", "rsa", "public").Return([]byte(nil), token.ErrKeyUnavailable)
	if _, err := x.Validate(""); err != token.ErrKeyUnavailable {
		t.Error(err)
	}
}

func TestValidateCorruptToken(t *testing.T) {
	m, _ := mock.New(hclog.NewNullLogger())
	x, err := NewRSA(hclog.NewNullLogger(), m)
	if err != nil {
		t.Fatal(err)
	}

	m.(*mock.Provider).On("Provide", "rsa", "public").Return([]byte(pubkey1), nil)
	if _, err := x.Validate(""); err != token.ErrInternalError {
		t.Error(err)
	}
}

func TestValidateWrongSigningMethod(t *testing.T) {
	m, _ := mock.New(hclog.NewNullLogger())
	x, err := NewRSA(hclog.NewNullLogger(), m)
	if err != nil {
		t.Fatal(err)
	}

	badToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	m.(*mock.Provider).On("Provide", "rsa", "public").Return([]byte(pubkey1), nil)
	if _, err := x.Validate(badToken); err != token.ErrTokenInvalid {
		t.Logf("%T", err)
		t.Error(err)
	}
}

func TestValidateExpiredToken(t *testing.T) {
	m, _ := mock.New(hclog.NewNullLogger())
	x, err := NewRSA(hclog.NewNullLogger(), m)
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

	m.(*mock.Provider).On("Provide", "rsa", "private").Return([]byte(privkey1), nil)
	tkn, err := x.Generate(c, cfg)
	if err != nil {
		t.Fatal(err)
	}

	m.(*mock.Provider).On("Provide", "rsa", "public").Return([]byte(pubkey1), nil)
	if _, err := x.Validate(tkn); err != token.ErrTokenInvalid {
		t.Fatalf("Incorrect error returned, expected '%v', got '%v'", token.ErrTokenInvalid, err)
	}
}

func TestPubKeyNoSuchKey(t *testing.T) {
	m, _ := mock.New(hclog.NewNullLogger())
	x, err := NewRSA(hclog.NewNullLogger(), m)
	if err != nil {
		t.Fatal(err)
	}

	m.(*mock.Provider).On("Provide", "rsa", "public").Return([]byte(nil), keyprovider.ErrNoSuchKey)
	k, err := x.(*RSATokenService).pubkey()
	if k != nil || err != token.ErrKeyUnavailable {
		t.Fatal(err)
	}
}

func TestPubKeyBadDecode(t *testing.T) {
	m, _ := mock.New(hclog.NewNullLogger())
	x, err := NewRSA(hclog.NewNullLogger(), m)
	if err != nil {
		t.Fatal(err)
	}

	m.(*mock.Provider).On("Provide", "rsa", "public").Return([]byte("for sure not a PEM block"), nil)
	k, err := x.(*RSATokenService).pubkey()
	if k != nil || err != token.ErrKeyUnavailable {
		t.Fatal(err)
	}
}

func TestPubKeyBadParse(t *testing.T) {
	m, _ := mock.New(hclog.NewNullLogger())
	x, err := NewRSA(hclog.NewNullLogger(), m)
	if err != nil {
		t.Fatal(err)
	}

	m.(*mock.Provider).On("Provide", "rsa", "public").Return([]byte(pubkey3), nil)
	k, err := x.(*RSATokenService).pubkey()
	if k != nil || err != token.ErrKeyUnavailable {
		t.Fatal(err)
	}
}

func TestPubKeyBadKeyType(t *testing.T) {
	m, _ := mock.New(hclog.NewNullLogger())
	x, err := NewRSA(hclog.NewNullLogger(), m)
	if err != nil {
		t.Fatal(err)
	}

	m.(*mock.Provider).On("Provide", "rsa", "public").Return([]byte(pubkey2), nil)
	k, err := x.(*RSATokenService).pubkey()
	if k != nil || err != token.ErrKeyUnavailable {
		t.Fatal(err)
	}
}

func TestPrivKeyUnknownError(t *testing.T) {
	m, _ := mock.New(hclog.NewNullLogger())
	x, err := NewRSA(hclog.NewNullLogger(), m)
	if err != nil {
		t.Fatal(err)
	}

	m.(*mock.Provider).On("Provide", "rsa", "private").Return([]byte(nil), errors.New("wat"))
	k, err := x.(*RSATokenService).privkey()
	if k != nil || err.Error() != "wat" {
		t.Fatal(err)
	}
}

func TestPrivKeyBadDecode(t *testing.T) {
	m, _ := mock.New(hclog.NewNullLogger())
	x, err := NewRSA(hclog.NewNullLogger(), m)
	if err != nil {
		t.Fatal(err)
	}

	m.(*mock.Provider).On("Provide", "rsa", "private").Return([]byte("for sure not a PEM block"), nil)
	k, err := x.(*RSATokenService).privkey()
	if k != nil || err != token.ErrKeyUnavailable {
		t.Fatal(err)
	}
}
