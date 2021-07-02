package bitcask

import (
	"github.com/hashicorp/go-hclog"

	"github.com/netauth/netauth/internal/startup"
)

func init() {
	startup.RegisterCallback(cb)
}

func cb() {
	hclog.L().Info("Bitcask is not supported in this environment")
}
