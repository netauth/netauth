package main

import (
	"os"

	"github.com/hashicorp/go-hclog"
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
			Name:  "nsutil",
			Level: hclog.LevelFromString(level),
		})
	}
	hclog.SetDefault(appLogger)
	appLogger.Debug("Build information as follows", "version", version, "commit", commit, "builddate", date)

	execute()
}
