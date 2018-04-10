package jwt

import (
	"crypto/rsa"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/NetAuth/NetAuth/internal/token"
)

var (
	config = token.TokenConfig{
		Lifetime: time.Minute * 5,
		Renewals: 0,
		Issuer:   "NetAuth Test",
	}

	testToken = "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJFbnRpdHlJRCI6ImZvbyIsIkNhcGFiaWxpdGllcyI6bnVsbCwiUmVuZXdhbHNMZWZ0IjowLCJhdWQiOiJVbnJlc3RyaWN0ZWQiLCJleHAiOi02MjEzNTU5NjUwMCwianRpIjoiZm9vIiwiaWF0IjotNjIxMzU1OTY4MDAsImlzcyI6Ik5ldEF1dGggVGVzdCIsIm5iZiI6LTYyMTM1NTk2ODAwLCJzdWIiOiJOZXRBdXRoIFN0YW5kYXJkIFRva2VuIn0.OpGakOumqqA9EscEU3vgDkX3DJtVifLxmpLXgPr5YZ7bgWxXk-pWBxSG4aAgbdSC2G78JGi6QuJXc849XvtuWqdDZI8pTAWnNZSnicdJr0cHdGCnvgOe4Iwj2U6TAgAfwYAXe_JZJM8HRQXHULUihGIyQSkJgqrIlVoJCidYXoaTThUplWYqvWpaim6LmujC2ko3oJq7bCDzi1FuMiGrTwedHRKiFBBJet3tGsQUXLfhVR9qWz44iRyAaRCyMcTkjEN3tMPBXVYBy1ms_b8ZaQvPKWnJzP9EHjUfIO2u0hmQUWUfoc0ZDqbK0uXUOgNCrwYxolHD2U1c71luA3tDxQ"
)

func init() {
	// We're throwing the error away here since this is parsing
	// the reference format, if that doesn't work its very
	// unlikely anything else does.
	t, _ := time.Parse(time.ANSIC, time.ANSIC)

	config.IssuedAt = t
	config.NotBefore = t
}

func mkTmpTestDir(t *testing.T) string {
	dir, err := ioutil.TempDir("/tmp", "tkntest")
	if err != nil {
		t.Error(err)
	}
	return dir
}

func cleanTmpTestDir(dir string, t *testing.T) {
	// Remove the tmpdir, don't want to clutter the filesystem
	if err := os.RemoveAll(dir); err != nil {
		t.Log(err)
	}
}

func genFixedKey(t *testing.T) {
	// Chosen by fair dice role.
	r := rand.New(rand.NewSource(4))

	key, err := rsa.GenerateKey(r, 2048)
	if err != nil {
		t.Error(err)
	}

	d, err := json.Marshal(key)
	if err != nil {
		t.Error(err)
	}
	if err := ioutil.WriteFile(*key_blob, d, 0400); err != nil {
		t.Error(err)
	}
}

func TestNewGenerateKeys(t *testing.T) {
	testDir := mkTmpTestDir(t)
	defer cleanTmpTestDir(testDir, t)
	*key_blob = testDir + "/key.dat"
	*generate = true

	_, err := NewRSA()
	if err != nil {
		t.Error(err)
	}
}

func TestNewMissingKeys(t *testing.T) {
	testDir := mkTmpTestDir(t)
	defer cleanTmpTestDir(testDir, t)
	*key_blob = testDir + "/key.dat"
	*generate = false

	_, err := NewRSA()
	if err != token.NO_GENERATE_KEYS {
		t.Error(err)
	}
}

func TestLoadExistingKey(t *testing.T) {
	testDir := mkTmpTestDir(t)
	defer cleanTmpTestDir(testDir, t)
	*key_blob = testDir + "/key.dat"

	// This one should generate keys
	*generate = true
	_, err := NewRSA()
	if err != nil {
		t.Error(err)
	}

	// This one should be loading the existing key
	*generate = false
	_, err = NewRSA()
	if err != nil {
		t.Error(err)
	}
}

func TestGenerateToken(t *testing.T) {
	testDir := mkTmpTestDir(t)
	defer cleanTmpTestDir(testDir, t)
	*key_blob = testDir + "/key.dat"

	// generate a fixed value key
	genFixedKey(t)

	// Create the token service which will use the key generated
	// earlier
	x, err := NewRSA()
	if err != nil {
		t.Error(err)
	}

	c := token.Claims{
		EntityID: "foo",
	}
	tkn, err := x.Generate(c, config)
	if err != nil {
		t.Error(err)
	}

	if tkn != testToken {
		t.Error("Bad Token")
	}
}

func TestValidateKey(t *testing.T) {
	testDir := mkTmpTestDir(t)
	defer cleanTmpTestDir(testDir, t)
	*key_blob = testDir + "/key.dat"

	// generate a fixed value key
	genFixedKey(t)

	// Create the token service which will use the key generated
	// earlier
	x, err := NewRSA()
	if err != nil {
		t.Error(err)
	}

	c := token.Claims{
		EntityID: "foo",
	}

	cfg := token.TokenConfig{
		Lifetime:  time.Minute * 5,
		IssuedAt:  time.Now(),
		NotBefore: time.Now(),
		Renewals:  0,
		Issuer:    "NetAuth Test",
	}

	tkn, err := x.Generate(c, cfg)
	if err != nil {
		t.Error(err)
	}

	claims, err := x.Validate(tkn)
	if err != nil {
		t.Error(err)
	}

	// structs containing []string can't be compared directly, so
	// we compare the value that was set earlier
	if claims.EntityID != c.EntityID {
		t.Error("Claims are not the same!")
	}
}
