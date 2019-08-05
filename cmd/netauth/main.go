package main

import (
	"fmt"
	"os"
	"github.com/NetAuth/NetAuth/internal/ctl"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if _, set := os.LookupEnv("NETAUTH_VERBOSE"); set {
		fmt.Printf("NetAuth %v:%v Built on %v", version, commit, date)
	}
	ctl.Execute()
}
