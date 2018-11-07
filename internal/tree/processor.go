package tree

import (
	"log"
	"sort"

	"github.com/NetAuth/NetAuth/internal/tree/errors"

	pb "github.com/NetAuth/Protocol"
)

// FetchHooks configures an EntityProcessor with hook chains from an
// external map.
func (ep *EntityProcessor) FetchHooks(chain string, hookmap map[string][]EntityProcessorHook) error {
	hookChain, ok := hookmap[chain]
	if !ok {
		return tree.ErrUnknownHookChain
	}

	if len(hookChain) == 0 {
		return tree.ErrEmptyHookChain
	}

	ep.hooks = hookChain
	return nil
}

// Run handles entity processor pipelines
func (ep *EntityProcessor) Run() (*pb.Entity, error) {
	for _, h := range ep.hooks {
		//log.Println(h.Name(), ep.Entity)
		if err := h.Run(ep.Entity, ep.RequestData); err != nil {
			return nil, err
		}
		//log.Println(h.Name(), ep.Entity)
	}
	return ep.Entity, nil
}

// Register adds a new hook to the processing pipeline
func (ep *EntityProcessor) Register(h EntityProcessorHook) error {
	m := make(map[string]bool)
	for _, rh := range ep.hooks {
		m[rh.Name()] = true
	}

	if _, ok := m[h.Name()]; ok {
		// Already registered, can't have two of the same hook
		return tree.ErrHookExists
	}

	ep.hooks = append(ep.hooks, h)

	sort.Slice(ep.hooks, func(i, j int) bool {
		return ep.hooks[i].Priority() < ep.hooks[j].Priority()
	})

	return nil
}

// EntityHookMustRegister registers hooks to named chains, and is
// intended to be used during startup to register changes or abort the
// service process.
func (m *Manager) EntityHookMustRegister(chain string, hook EntityProcessorHook) {
	mp := make(map[string]bool)
	for _, rh := range m.entityProcesses[chain] {
		mp[rh.Name()] = true
	}

	if _, ok := mp[hook.Name()]; ok {
		// Already registered, can't have two of the same hook
		log.Fatalf("Hook %s already exists in chain %s", hook.Name(), chain)
	}

	m.entityProcesses[chain] = append(m.entityProcesses[chain], hook)

	sort.Slice(m.entityProcesses[chain], func(i, j int) bool {
		return m.entityProcesses[chain][i].Priority() < m.entityProcesses[chain][j].Priority()
	})
}
