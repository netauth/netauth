package consumer

import (
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/NetAuth/NetAuth/internal/health"
	"github.com/NetAuth/NetAuth/internal/plugin/tree/common"
)

var (
	handshakeConfig = plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "TREE_PLUGIN",
		MagicCookieValue: "treehello",
	}

	pluginMap = map[string]plugin.Plugin{
		"treep": &common.GoPluginRPC{},
	}
)

// New returns a new plugin reference.
func New(path string) (Ref, error) {
	return Ref{
		path: path,
		log:  hclog.L().Named("plugin"),
	}, nil
}

// Shutdown closes down the plugin which prevents us from leaking
// processes.
func (r *Ref) Shutdown() {
	r.cfg.Kill()
}

// Init sets up the plugin and gets it running, an expectation exists
// that the ref will be stored during the lifetime of the server.
func (r *Ref) Init() error {
	health.RegisterCheck("plugin-"+r.Name(), r.healthCheck)

	r.cfg = plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command(r.path),
		Logger:          r.log,
	})

	var err error
	r.client, err = r.cfg.Client()
	if err != nil {
		return err
	}

	raw, err := r.client.Dispense("treep")
	if err != nil {
		return err
	}

	impl := raw.(*common.GoPluginClient)
	r.GoPluginClient = impl

	return nil
}

// Name returns a usable name for the plugin
func (r *Ref) Name() string {
	base := filepath.Base(r.path)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

func (r *Ref) healthCheck() health.SubsystemStatus {
	status := health.SubsystemStatus{
		Name:   "plugin-" + r.Name(),
		OK:     true,
		Status: "Plugin is OK",
	}

	err := r.client.Ping()
	if err != nil {
		status.OK = false
		status.Status = err.Error()
	}

	return status
}
