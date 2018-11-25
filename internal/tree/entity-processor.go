package tree

import (
	"log"
	"sort"

	pb "github.com/NetAuth/Protocol"
)

// EntityHookConstructor functions construct EntityHook instances and
// return the hooks for registration into the map of hooks.  This
// allows the hooks to notify the module of thier presence and defer
// construction until a RefContext can be prepared.
type EntityHookConstructor func(RefContext) (EntityHook, error)

// An EntityHook is a function that transforms an entity as
// part of an EntityProcessor pipeline.
type EntityHook interface {
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
	if *debugChains {
		log.Println("Executing EntityHookConstructor callbacks")
	}
	for _, v := range eHookConstructors {
		hook, err := v(m.refContext)
		if err != nil {
			log.Println(err)
			continue
		}
		m.entityProcessorHooks[hook.Name()] = hook
	}
	if *debugChains {
		log.Printf("The following (entity) hooks are loaded:")
		for name := range m.entityProcessorHooks {
			log.Printf("  %s", name)
		}
	}
}

// InitializeEntityChains initializes the map of chains stored on the
// manager.  It is expected that any merging of an external
// configuration has happened before this function is called.
func (m *Manager) InitializeEntityChains(c ChainConfig) error {
	for chain, hooks := range c {
		if *debugChains {
			log.Printf("Initializing chain '%s'", chain)

		}
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
		if *debugChains {
			for _, hook := range m.entityProcesses[chain] {
				log.Printf("  %s", hook.Name())
			}
		}
	}
	return nil
}

// CheckRequiredEntityChains searches for all chains in the default
// chains list and logs a fatal error if one isn't found in the
// configured list.  This allows the system to later assert the
// presence of chains without checking, since they cannot be modified
// after loading.
func (m *Manager) CheckRequiredEntityChains() {
	for k := range defaultEntityChains {
		if _, ok := m.entityProcesses[k]; !ok {
			log.Fatalf("Required chain %s is not loaded", k)
		}
		if len(m.entityProcesses[k]) == 0 {
			log.Fatalf("Required chain %s is empty", k)
		}
	}
}

// RunEntityChain runs the specified chain with de specifying values
// to be consumed by the chain.
func (m *Manager) RunEntityChain(chain string, de *pb.Entity) (*pb.Entity, error) {
	e := new(pb.Entity)
	hookChain := m.entityProcesses[chain]
	for _, h := range hookChain {
		if err := h.Run(e, de); err != nil {
			return nil, err
		}
	}
	return e, nil
}
