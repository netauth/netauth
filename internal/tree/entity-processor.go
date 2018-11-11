package tree

import (
	"log"
	"sort"

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

// RegisterEntityHookConstructor registers the entity hook
// constructors to be called during the initialization of the main
// tree manager.
func RegisterEntityHookConstructor(name string, c EntityHookConstructor) {
	if _, ok := eHookConstructors[name]; ok {
		// Already registered
		log.Printf("A constructor for %s is already registered", name)
		return
	}
	eHookConstructors[name] = c
}

// InitializeEntityHooks runs all the EntityHookConstructors and
// registers the resulting hooks by name into m.entityProcessorHooks
func (m *Manager) InitializeEntityHooks() {
	log.Println("Executing EntityHookConstructor callbacks")
	for _, v := range eHookConstructors {
		hook, err := v(m.refContext)
		if err != nil {
			log.Println(err)
			continue
		}
		m.entityProcessorHooks[hook.Name()] = hook
	}
	log.Printf("The following (entity) hooks are loaded:")
	for name := range m.entityProcessorHooks {
		log.Printf("  %s", name)
	}
}

// InitializeEntityChains initializes the map of chains stored on the
// manager.  It is expected that any merging of an external
// configuration has happened before this function is called.
func (m *Manager) InitializeEntityChains(c ChainConfig) error {
	for chain, hooks := range c {
		log.Printf("Initializing chain '%s'", chain)
		for _, h := range hooks {
			eph, ok := m.entityProcessorHooks[h]
			if !ok {
				log.Printf("There is no hook named '%s'", h)
				return ErrUnknownHook
			}
			m.entityProcesses[chain] = append(m.entityProcesses[chain], eph)
		}
		sort.Slice(m.entityProcesses[chain], func(i, j int) bool {
			return m.entityProcesses[chain][i].Priority() < m.entityProcesses[chain][j].Priority()
		})
		for _, hook := range m.entityProcesses[chain] {
			log.Printf("  %s", hook.Name())
		}
	}
	return nil
}

// FetchHooks configures an EntityProcessor with hook chains from an
// external map.
func (ep *EntityProcessor) FetchHooks(chain string, hookmap map[string][]EntityProcessorHook) error {
	hookChain, ok := hookmap[chain]
	if !ok {
		return ErrUnknownHookChain
	}

	if len(hookChain) == 0 {
		return ErrEmptyHookChain
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
