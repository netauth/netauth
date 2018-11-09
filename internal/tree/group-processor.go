package tree

import (
	"log"
	"sort"

	"github.com/NetAuth/NetAuth/internal/tree/errors"

	pb "github.com/NetAuth/Protocol"
)

// A GroupProcessor is a chain of functions that performs mutations on
// a group.
type GroupProcessor struct {
	Group       *pb.Group
	RequestData *pb.Group
	hooks       []GroupProcessorHook
}

// A GroupProcessorHook is a function that transforms a group as part
// of a GroupProcessor Pipeline.
type GroupProcessorHook interface {
	Priority() int
	Name() string
	Run(*pb.Group, *pb.Group) error
}

// FetchHooks configures an GroupProcessor with hook chains from an
// external map.
func (ep *GroupProcessor) FetchHooks(chain string, hookmap map[string][]GroupProcessorHook) error {
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
func (ep *GroupProcessor) Run() (*pb.Group, error) {
	for _, h := range ep.hooks {
		//log.Println(h.Name(), ep.Group)
		if err := h.Run(ep.Group, ep.RequestData); err != nil {
			return nil, err
		}
		//log.Println(h.Name(), ep.Group)
	}
	return ep.Group, nil
}

// Register adds a new hook to the processing pipeline
func (ep *GroupProcessor) Register(h GroupProcessorHook) error {
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

// GroupHookMustRegister registers hooks to named chains, and is
// intended to be used during startup to register changes or abort the
// service process.
func (m *Manager) GroupHookMustRegister(chain string, hook GroupProcessorHook) {
	mp := make(map[string]bool)
	for _, rh := range m.groupProcesses[chain] {
		mp[rh.Name()] = true
	}

	if _, ok := mp[hook.Name()]; ok {
		// Already registered, can't have two of the same hook
		log.Fatalf("Hook %s already exists in chain %s", hook.Name(), chain)
	}

	m.groupProcesses[chain] = append(m.groupProcesses[chain], hook)

	sort.Slice(m.entityProcesses[chain], func(i, j int) bool {
		return m.entityProcesses[chain][i].Priority() < m.entityProcesses[chain][j].Priority()
	})
}
