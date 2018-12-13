package util

import (
	"testing"
)

func TestNextEntityNumber(t *testing.T) {
	if n, err := NextEntityNumber(okEntityLoader, okEntityIDFinder); err != nil || n != 3 {
		t.Error("Bad result from NextEntityNumber")
	}
}

func TestNextEntityNumberBadIDs(t *testing.T) {
	if _, err := NextEntityNumber(okEntityLoader, evilEntityIDFinder); err == nil {
		t.Error("Got nil error from evil ID finder")
	}
}

func TestNextEntityNumberBadLoader(t *testing.T) {
	if _, err := NextEntityNumber(evilEntityLoader, okEntityIDFinder); err == nil {
		t.Error("Got nil error from evil loader")
	}
}

func TestNextGroupNumber(t *testing.T) {
	if n, err := NextGroupNumber(okGroupLoader, okGroupIDFinder); err != nil || n != 3 {
		t.Error("Bad result from NextGroupNumber")
	}
}

func TestNextGroupNumberBadIDs(t *testing.T) {
	if _, err := NextGroupNumber(okGroupLoader, evilGroupIDFinder); err == nil {
		t.Error("Got nil error from evil ID finder")
	}
}

func TestNextGroupNumberBadLoader(t *testing.T) {
	if _, err := NextGroupNumber(evilGroupLoader, okGroupIDFinder); err == nil {
		t.Error("Got nil error from evil loader")
	}
}
