package main

import (
	"os"

	"github.com/hashicorp/go-hclog"

	"github.com/NetAuth/NetAuth/internal/ctl"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	level, set := os.LookupEnv("NETAUTH_VERBOSE")
	if !set {
		hclog.SetDefault(hclog.NewNullLogger())
	} else {
		appLogger := hclog.New(&hclog.LoggerOptions{
			Name:  "netauth",
			Level: hclog.LevelFromString(level),
		})
		hclog.SetDefault(appLogger)
	}

	if _, set := os.LookupEnv("NETAUTH_VERBOSE"); set {
		hclog.L().Debug("Build information as follows", "version", version, "commit", commit, "builddate", date)
	}
	ctl.Execute()
}
