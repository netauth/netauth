package db

import (
	"github.com/hashicorp/go-hclog"
)

var (
	kvBackends map[string]KVFactory
)

func init() {
	kvBackends = make(map[string]KVFactory)
}

// RegisterKV registers a KV factory which can be called later.
func RegisterKV(name string, factory KVFactory) {
	if _, ok := kvBackends[name]; ok {
		return
	}
	log().Info("Registered KV Store", "kv", name)
	kvBackends[name] = factory
}

func newKV(name string, l hclog.Logger) (KVStore, error) {
	f, ok := kvBackends[name]
	if !ok {
		log().Debug("Requested bad backend", "backend", name, "known", kvBackends)
		return nil, ErrUnknownDatabase
	}
	log().Debug("Initializing database with backend", "backend", name)
	return f(log())
}
