package main

import (
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/proto"

	"github.com/netauth/netauth/internal/tree/util"
	"github.com/netauth/netauth/pkg/plugin/tree"

	pb "github.com/netauth/protocol"
)

var (
	appLogger hclog.Logger

	cfg *viper.Viper
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.netauth")
	viper.AddConfigPath("/etc/netauth/")

	if err := viper.ReadInConfig(); err != nil {
		appLogger.Error("Fatal error reading configuration", "error", err)
		os.Exit(1)
	}

	viper.SetDefault("log.level", "INFO")
	appLogger = hclog.New(&hclog.LoggerOptions{
		Name:  "fail2lock",
		Level: hclog.LevelFromString(viper.GetString("log.level")),
	})
	hclog.SetDefault(appLogger)

	viper.SetDefault("plugin.fail2lock.allowed_fails", 3)
	viper.SetDefault("plugin.fail2lock.interval", time.Minute*15)
	cfg = viper.Sub("plugin.fail2lock")
}

func main() {
	appLogger.Info("fail2lock initialized",
		"allowed_fails", cfg.GetInt("allowed_fails"),
		"interval", cfg.GetDuration("interval"))

	tree.PluginMain(fail2lock{
		NullPlugin: tree.NullPlugin{},
		failDB:     make(map[string][]time.Time),
	})
}

type fail2lock struct {
	tree.NullPlugin

	failDB map[string][]time.Time
}

// PreAuthCheck runs first and checks for a number of failed markers
// for auth.  Assuming it finds more than the configured amount, it
// will lock the entity and return which will cause an authentication
// failure.
func (f fail2lock) PreAuthCheck(e, de pb.Entity) (pb.Entity, error) {
	flags := util.PatchKeyValueSlice(e.Meta.UntypedMeta, "READ", "fail2lock", "")

	if len(flags) == 1 && strings.Split(flags[0], ":")[1] == "RESET" {
		delete(f.failDB, e.GetID())
		e.Meta.UntypedMeta = util.PatchKeyValueSlice(e.Meta.UntypedMeta, "CLEARFUZZY", "fail2lock", "")
		return e, nil
	}

	inIntervalFails := 0
	startTime := time.Now().Add(cfg.GetDuration("interval") * -1)
	for _, t := range f.failDB[e.GetID()] {
		if t.After(startTime) {
			inIntervalFails++
		}
	}

	if inIntervalFails >= cfg.GetInt("allowed_fails") {
		appLogger.Warn("fail2lock is locking an entity",
			"entity", e.GetID(),
			"fails", inIntervalFails,
			"allowed", cfg.GetInt("allowed_fails"))
		e.Meta.Locked = proto.Bool(true)
		return e, nil
	}

	f.failDB[e.GetID()] = append(f.failDB[e.GetID()], time.Now())

	return e, nil
}

// PostAuthCheck is only run if the authentication has succeeded.  In
// this case we will remove the fail2lock key from the entity's
// metadata since they have successfully authenticated.
func (f fail2lock) PostAuthCheck(e, de pb.Entity) (pb.Entity, error) {
	delete(f.failDB, e.GetID())
	return e, nil
}
