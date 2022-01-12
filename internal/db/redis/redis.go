package redis

import (
	"context"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-redis/redis"
	"github.com/hashicorp/go-hclog"
	"github.com/netauth/netauth/internal/db"
	"github.com/netauth/netauth/internal/startup"
	"github.com/spf13/viper"
)

func init() {
	startup.RegisterCallback(cb)
	db.RegisterKV("redis", New)
}

func cb() {
	db.RegisterKV("redis", New)
}

type RedisStore struct {
	redis *redis.Client
	l     hclog.Logger
	eF    func(db.Event)
}

type eventType int

const (
	eventUpdate eventType = iota
	eventDelete
)

func New(l hclog.Logger) (db.KVStore, error) {
	x := &RedisStore{l: l.Named("redis")}
	url := viper.GetViper().GetString("redis.url")
	opt, err := redis.ParseURL(url)
	x.redis = redis.NewClient(opt)

	return x, err
}

func (r *RedisStore) Put(_ context.Context, k string, v []byte) error {
	if err := r.redis.Set(k, v, 0).Err(); err != nil {
		return err
	}

	r.fireEventForKey(k, eventUpdate)
	return nil
}

func (r *RedisStore) Get(_ context.Context, k string) ([]byte, error) {
	if v := r.redis.Get(k); v.Err() == nil {
		return v.Bytes()
	}

	return nil, db.ErrNoValue
}

func (r *RedisStore) Del(_ context.Context, k string) error {
	if err := r.redis.Del(k).Err(); err != nil {
		return err
	}
	r.fireEventForKey(k, eventDelete)
	return nil
}

func (r *RedisStore) Keys(_ context.Context, f string) ([]string, error) {
	keys, err := r.redis.Keys("*").Result()
	out := []string{}

	for _, k := range keys {
		if m, _ := path.Match(f, k); m {
			out = append(out, k)
		}
	}

	return out, err
}

func (r *RedisStore) Close() error { return r.redis.Close() }

// Capabilities returns that this key/value store supports te mutable
// property, allowing it to be writeable to the higher level systems.
func (r *RedisStore) Capabilities() []db.KVCapability {
	return []db.KVCapability{db.KVMutable}
}

// fireEventForKey maps from a key to an entity or group and fires an
// appropriate event for the given key.
func (r *RedisStore) fireEventForKey(k string, t eventType) {
	switch {
	case t == eventUpdate && strings.HasPrefix(k, "/entities/"):
		r.eF(db.Event{
			PK:   filepath.Base(k),
			Type: db.EventEntityUpdate,
		})
	case t == eventDelete && strings.HasPrefix(k, "/entities/"):
		r.eF(db.Event{
			PK:   filepath.Base(k),
			Type: db.EventEntityDestroy,
		})
	case t == eventUpdate && strings.HasPrefix(k, "/groups/"):
		r.eF(db.Event{
			PK:   filepath.Base(k),
			Type: db.EventGroupUpdate,
		})
	case t == eventDelete && strings.HasPrefix(k, "/groups/"):
		r.eF(db.Event{
			PK:   filepath.Base(k),
			Type: db.EventGroupDestroy,
		})
	default:
		r.l.Warn("Event translation called with unknown key prefix", "type", t, "key", k)
	}
}

// SetEventFunc sets up a function to call to fire events to
// subscribers.
func (r *RedisStore) SetEventFunc(f func(db.Event)) {
	r.eF = f
}
