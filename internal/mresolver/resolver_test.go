package mresolver

import (
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

func TestSetParentLogger(t *testing.T) {
	x := New()

	x.SetParentLogger(hclog.L())
	assert.Equal(t, hclog.L().Named("resolver"), x.l)
}
