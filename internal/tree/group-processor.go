package tree

import (
	"log"
	"sort"

	pb "github.com/NetAuth/Protocol"
)

// GroupHookConstructor functions construct GroupProcessorHook
// instances and return the hooks for registration into the map of
// hooks.  This allows the hooks to notify the module of thier
// presence and defer construction until a RefContext can be prepared.
type GroupHookConstructor func(RefContext) (GroupProcessorHook, error)

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
	log.Println("Executing GroupHookConstructor callbacks")
	for _, v := range gHookConstructors {
		hook, err := v(m.refContext)
		if err != nil {
			log.Println(err)
			continue
		}
		m.groupProcessorHooks[hook.Name()] = hook
	}
	log.Printf("The following (group) hooks are loaded:")
	for name := range m.groupProcessorHooks {
		log.Printf("  %s", name)
	}
}

// InitializeGroupChains initializes the map of chains stored on the
// manager.  It is expected that any merging of an external
// configuration has happened before this function is called.
func (m *Manager) InitializeGroupChains(c ChainConfig) error {
	for chain, hooks := range c {
		log.Printf("Initializing chain '%s'", chain)
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
		for _, hook := range m.groupProcesses[chain] {
			log.Printf("  %s", hook.Name())
		}
	}
	return nil
}

// FetchHooks configures an GroupProcessor with hook chains from an
// external map.
func (ep *GroupProcessor) FetchHooks(chain string, hookmap map[string][]GroupProcessorHook) error {
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
