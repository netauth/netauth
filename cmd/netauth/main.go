package main

import (
	"os"

	"github.com/hashicorp/go-hclog"

	"github.com/netauth/netauth/internal/ctl"

	_ "github.com/netauth/netauth/pkg/token/cache/fs"
	_ "github.com/netauth/netauth/pkg/token/jwt"
	_ "github.com/netauth/netauth/pkg/token/keyprovider/fs"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"

	appLogger hclog.Logger
)

func main() {
	level, set := os.LookupEnv("NETAUTH_LOGLEVEL")
	if !set {
		appLogger = hclog.NewNullLogger()
	} else {
		appLogger = hclog.New(&hclog.LoggerOptions{
			Name:  "netauth",
			Level: hclog.LevelFromString(level),
		})
	}
	hclog.SetDefault(appLogger)
	appLogger.Debug("Build information as follows", "version", version, "commit", commit, "builddate", date)

	ctl.Execute()
}
