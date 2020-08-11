package tree

import (
	"testing"

	"github.com/hashicorp/go-hclog"
)

func TestSetParentLogger(t *testing.T) {
	initlb = nil

	l := hclog.NewNullLogger()
	SetParentLogger(l)
	if log() == nil {
		t.Error("log was not set")
	}
}

func TestLogParentUnset(t *testing.T) {
	initlb = nil

	if log() == nil {
		t.Error("auto log was not aquired")
	}
}
