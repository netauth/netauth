package tree

import (
	"log"
	"sort"

	"github.com/NetAuth/NetAuth/internal/tree/errors"

	pb "github.com/NetAuth/Protocol"
)

// EntityHookConstructor functions construct EntityProcessorHook
// instances and return the hooks for registration into the map of
// hooks.  This allows the hooks to notify the module of thier
// presence and defer construction until a RefContext can be prepared.
type EntityHookConstructor func(RefContext) (EntityProcessorHook, error)

// An EntityProcessor is a chain of functions that modify entities in
// some way.
type EntityProcessor struct {
	Entity      *pb.Entity
	RequestData *pb.Entity
	hooks       []EntityProcessorHook
}

// An EntityProcessorHook is a function that transforms an entity as
// part of an EntityProcessor pipeline.
type EntityProcessorHook interface {
	Priority() int
	Name() string
	Run(*pb.Entity, *pb.Entity) error
}

var (
	eHookConstructors map[string]EntityHookConstructor
)

func init() {
	eHookConstructors = make(map[string]EntityHookConstructor)
}

// InitializeEntityHooks runs all the EntityHookConstructors and
// registers the resulting hooks by name into m.entityProcessorHooks
func (m *Manager) InitializeEntityHooks() {
}

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

