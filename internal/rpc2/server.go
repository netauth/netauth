package rpc2

import (
	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"
)

// New returns a ready to use server implementation.
func New(r Refs) *Server {
	return &Server{
		Service:  r.TokenService,
		Manager:  r.Tree,
		readonly: viper.GetBool("server.readonly"),
		log:      hclog.L().Named("rpc2"),
	}
}
