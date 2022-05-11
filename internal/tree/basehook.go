package tree

import (
	"github.com/hashicorp/go-hclog"

	"github.com/netauth/netauth/internal/crypto"
)

// The BaseHook contains the critical fields needed to register and
// run hook pipelines.
type BaseHook struct {
	name     string
	priority int
	log      hclog.Logger
	storage  DB
	crypto   crypto.EMCrypto
}

// Name returns the name of a hook.  Names should be kabob case.
func (h *BaseHook) Name() string { return h.name }

// Priority returns the priority of a hook.  Priorities are banded as follows:
// 0-10:
//   Loaders
// 11-19:
//   Load time integrity checks
// 20-29:
//   User defined pre processing
// 30-49:
//   Checks and data validation
// 50-89:
//   User defined post processing
// 90-99:
//   Serialization and storage
func (h *BaseHook) Priority() int { return h.priority }

func (h *BaseHook) Log() hclog.Logger { return h.log }

func (h *BaseHook) Storage() DB { return h.storage }

func (h *BaseHook) Crypto() crypto.EMCrypto { return h.crypto }

// NewBaseHook returns a BaseHook struct for compact initialization
// during callback constructors.
func NewBaseHook(opts ...HookOption) BaseHook {
	b := BaseHook{
		name:     "INVALID",
		priority: -1,
		log:      hclog.NewNullLogger(),
	}

	for _, o := range opts {
		o(&b)
	}

	return b
}

func WithHookLogger(l hclog.Logger) HookOption { return func(b *BaseHook) { b.log = l } }

func WithHookName(n string) HookOption { return func(b *BaseHook) { b.name = n } }

func WithHookPriority(p int) HookOption { return func(b *BaseHook) { b.priority = p } }

func WithHookStorage(d DB) HookOption { return func(b *BaseHook) { b.storage = d } }

func WithHookCrypto(c crypto.EMCrypto) HookOption { return func(b *BaseHook) { b.crypto = c } }
