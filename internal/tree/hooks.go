package tree

type BaseHook struct {
	name     string
	priority int
}

func (h *BaseHook) Name() string  { return h.name }
func (h *BaseHook) Priority() int { return h.priority }

func NewBaseHook(n string, p int) BaseHook {
	return BaseHook{
		name:     n,
		priority: p,
	}
}
