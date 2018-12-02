package tree

// The BaseHook contains the critical fields needed to register and
// run hook pipelines.
type BaseHook struct {
	name     string
	priority int
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

// NewBaseHook returns a BaseHook struct for compact initialization
// during callback constructors.
func NewBaseHook(n string, p int) BaseHook {
	return BaseHook{
		name:     n,
		priority: p,
	}
}
