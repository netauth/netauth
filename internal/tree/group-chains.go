package tree

import (
	"sort"

	pb "github.com/NetAuth/Protocol"
)

// GroupHookConstructor functions construct GroupHook instances and
// return the hooks for registration into the map of hooks.  This
// allows the hooks to notify the module of their presence and defer
// construction until a RefContext can be prepared.
type GroupHookConstructor func(RefContext) (GroupHook, error)

// An GroupHook is a function that transforms an group as
// part of an GroupProcessor pipeline.
type GroupHook interface {
	Priority() int
	Name() string
	Run(*pb.Group, *pb.Group) error
}

var (
	gHookConstructors map[string]GroupHookConstructor
)

func init() {
	gHookConstructors = make(map[string]GroupHookConstructor)
}

// RegisterGroupHookConstructor registers the group hook
// constructors to be called during the initialization of the main
// tree manager.
func RegisterGroupHookConstructor(name string, c GroupHookConstructor) {
	if _, ok := gHookConstructors[name]; ok {
		// Already registered
		logger.Trace("Duplicate GroupHookConstructor registration attempt", "hook", name)
		return
	}
	gHookConstructors[name] = c
	logger.Trace("GroupHookConstructor registered", "constructor", name)
}

// InitializeGroupHooks runs all the GroupHookConstructors and
// registers the resulting hooks by name into m.groupProcessorHooks
func (m *Manager) InitializeGroupHooks() {
	m.log.Debug("Executing GroupHookConstructor callbacks")
	for _, v := range gHookConstructors {
		hook, err := v(m.refContext)
		if err != nil {
			m.log.Warn("Error initializing hook", "hook", hook, "error", err)
			continue
		}
		m.groupHooks[hook.Name()] = hook
		m.log.Trace("GroupHook registered", "hook", hook.Name())
	}
}

// InitializeGroupChains initializes the map of chains stored on the
// manager.  It is expected that any merging of an external
// configuration has happened before this function is called.
func (m *Manager) InitializeGroupChains(c ChainConfig) error {
	for chain, hooks := range c {
		m.log.Debug("Initializing Group Chain", "chain", chain)
		for _, h := range hooks {
			if err := m.RegisterGroupHookToChain(h, chain); err != nil {
				return err
			}
		}
	}
	return nil
}

// RegisterGroupHookToChain registers a hook to a given chain.
func (m *Manager) RegisterGroupHookToChain(hook, chain string) error {
	eph, ok := m.groupHooks[hook]
	if !ok {
		m.log.Warn("Missing hook during chain initializtion", "chain", chain, "hook", hook)
		return ErrUnknownHook
	}
	m.groupProcesses[chain] = append(m.groupProcesses[chain], eph)
	sort.Slice(m.groupProcesses[chain], func(i, j int) bool {
		return m.groupProcesses[chain][i].Priority() < m.groupProcesses[chain][j].Priority()
	})
	m.log.Trace("Registered hook to chain", "chain", chain, "hook", hook)
	return nil
}

// CheckRequiredGroupChains searches for all chains in the default
// chains list and logs a fatal error if one isn't found in the
// configured list.  This allows the system to later assert the
// presence of chains without checking, since they cannot be modified
// after loading.
func (m *Manager) CheckRequiredGroupChains() error {
	for k := range defaultGroupChains {
		if _, ok := m.groupProcesses[k]; !ok {
			m.log.Error("Missing required chain", "chain", k)
			return ErrUnknownHookChain
		}
		if len(m.groupProcesses[k]) == 0 {
			m.log.Error("A required chain is empty", "chain", k)
			return ErrEmptyHookChain
		}
	}
	return nil
}

// RunGroupChain runs the specified chain with de specifying values
// to be consumed by the chain.
func (m *Manager) RunGroupChain(chain string, de *pb.Group) (*pb.Group, error) {
	e := new(pb.Group)
	hookChain := m.groupProcesses[chain]
	for _, h := range hookChain {
		m.log.Trace("Executing group hook", "chain", chain, "hook", h.Name())
		if err := h.Run(e, de); err != nil {
			m.log.Trace("Error during chain execution", "chain", chain, "hook", h.Name(), "error", err)
			return nil, err
		}
	}
	return e, nil
}
