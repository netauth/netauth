package tree

import (
	"log"
	"sort"

	pb "github.com/NetAuth/Protocol"
)

// GroupHookConstructor functions construct GroupHook
// instances and return the hooks for registration into the map of
// hooks.  This allows the hooks to notify the module of thier
// presence and defer construction until a RefContext can be prepared.
type GroupHookConstructor func(RefContext) (GroupHook, error)

// A GroupHook is a function that transforms a group as part
// of a Group Pipeline.
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

// RegisterGroupHookConstructor registers the entity hook
// constructors to be called during the initialization of the main
// tree manager.
func RegisterGroupHookConstructor(name string, c GroupHookConstructor) {
	if _, ok := gHookConstructors[name]; ok {
		// Already registered
		log.Printf("A constructor for %s is already registered", name)
		return
	}
	gHookConstructors[name] = c
}

// InitializeGroupHooks runs all the GroupHookConstructors and
// registers the resulting hooks by name into m.entityProcessorHooks
func (m *Manager) InitializeGroupHooks() {
	if *debugChains {
		log.Println("Executing GroupHookConstructor callbacks")
	}
	for _, v := range gHookConstructors {
		hook, err := v(m.refContext)
		if err != nil {
			log.Println(err)
			continue
		}
		m.groupProcessorHooks[hook.Name()] = hook
	}

	if *debugChains {
		log.Printf("The following (group) hooks are loaded:")
		for name := range m.groupProcessorHooks {
			log.Printf("  %s", name)
		}
	}
}

// InitializeGroupChains initializes the map of chains stored on the
// manager.  It is expected that any merging of an external
// configuration has happened before this function is called.
func (m *Manager) InitializeGroupChains(c ChainConfig) error {
	for chain, hooks := range c {
		if *debugChains {
			log.Printf("Initializing chain '%s'", chain)
		}
		for _, h := range hooks {
			gph, ok := m.groupProcessorHooks[h]
			if !ok {
				log.Printf("There is no hook named '%s'", h)
				return ErrUnknownHook
			}
			m.groupProcesses[chain] = append(m.groupProcesses[chain], gph)
		}
		sort.Slice(m.groupProcesses[chain], func(i, j int) bool {
			return m.groupProcesses[chain][i].Priority() < m.groupProcesses[chain][j].Priority()
		})
		if *debugChains {
			for _, hook := range m.groupProcesses[chain] {
				log.Printf("  %s", hook.Name())
			}
		}
	}
	return nil
}

// CheckRequiredGroupChains searches for all chains in the default
// chains list and logs a fatal error if one isn't found in the
// configured list.  This allows the system to later assert the
// presence of chains without checking, since they cannot be modified
// after loading.
func (m *Manager) CheckRequiredGroupChains() {
	for k := range defaultGroupChains {
		if _, ok := m.groupProcesses[k]; !ok {
			log.Fatalf("Required chain %s is not loaded", k)
		}
		if len(m.groupProcesses[k]) == 0 {
			log.Fatalf("Required chain %s is empty", k)
		}
	}
}

// RunGroupChain runs the specified chain with de specifying values
// to be consumed by the chain.
func (m *Manager) RunGroupChain(chain string, de *pb.Group) (*pb.Group, error) {
	e := new(pb.Group)
	hookChain := m.groupProcesses[chain]
	for _, h := range hookChain {
		if err := h.Run(e, de); err != nil {
			return nil, err
		}
	}
	return e, nil
}
