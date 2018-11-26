package tree

import (
	"errors"
	"testing"

	pb "github.com/NetAuth/Protocol"
)

func resetEntityConstructorMap() {
	eHookConstructors = make(map[string]EntityHookConstructor)
}

func TestECRegisterAndInitialize(t *testing.T) {
	*debugChains = true
	resetEntityConstructorMap()
	defer resetEntityConstructorMap()

	RegisterEntityHookConstructor("null-hook", goodEntityConstructor)
	RegisterEntityHookConstructor("null-hook", goodEntityConstructor)

	if len(eHookConstructors) != 1 {
		t.Error("Duplicate hook registered")
	}

	RegisterEntityHookConstructor("bad-hook", badEntityConstructor)

	if len(eHookConstructors) != 2 {
		t.Error("bad-hook wasn't registered")
	}

	em := Manager{
		entityHooks: make(map[string]EntityHook),
	}

	em.InitializeEntityHooks()
	if len(em.entityHooks) != 1 {
		t.Error("bad-hook was initialized")
	}
}

func TestECInitializeChainsOK(t *testing.T) {
	resetEntityConstructorMap()
	defer resetEntityConstructorMap()

	RegisterEntityHookConstructor("null-hook", goodEntityConstructor)
	RegisterEntityHookConstructor("null-hook2", goodEntityConstructor2)
	em := Manager{
		entityHooks:     make(map[string]EntityHook),
		entityProcesses: make(map[string][]EntityHook),
	}
	em.InitializeEntityHooks()

	c := map[string][]string{
		"TEST": []string{"null-hook", "null-hook2"},
	}

	if err := em.InitializeEntityChains(c); err != nil {
		t.Error(err)
	}
}

func TestECInitializeBadHook(t *testing.T) {
	resetEntityConstructorMap()
	defer resetEntityConstructorMap()

	em := Manager{
		entityHooks:     make(map[string]EntityHook),
		entityProcesses: make(map[string][]EntityHook),
	}
	em.InitializeEntityHooks()

	c := map[string][]string{
		"TEST": []string{"unknown-hook"},
	}

	if err := em.InitializeEntityChains(c); err != ErrUnknownHook {
		t.Error(err)
	}
}

func TestECCheckRequiredMissing(t *testing.T) {
	resetEntityConstructorMap()
	defer resetEntityConstructorMap()

	em := Manager{
		entityHooks:     make(map[string]EntityHook),
		entityProcesses: make(map[string][]EntityHook),
	}

	if err := em.CheckRequiredEntityChains(); err != ErrUnknownHookChain {
		t.Error("Passed with a required chain missing")
	}
}

func TestECCheckRequiredEmpty(t *testing.T) {
	resetEntityConstructorMap()
	defer resetEntityConstructorMap()

	em := Manager{
		entityHooks:     make(map[string]EntityHook),
		entityProcesses: make(map[string][]EntityHook),
	}

	// This lets us do this without having hooks loaded, we just
	// register something into all the chains, and then kill one
	// of them at the end.
	for k := range defaultEntityChains {
		em.entityProcesses[k] = []EntityHook{
			&nullEntityHook{},
		}
	}

	em.entityProcesses["CREATE"] = nil

	if err := em.CheckRequiredEntityChains(); err != ErrEmptyHookChain {
		t.Error("Passed with an empty required chain")
	}
}

type nullEntityHook struct{}

func (*nullEntityHook) Name() string              { return "null-hook" }
func (*nullEntityHook) Priority() int             { return 50 }
func (*nullEntityHook) Run(_, _ *pb.Entity) error { return nil }
func goodEntityConstructor(_ RefContext) (EntityHook, error) {
	return &nullEntityHook{}, nil
}

type nullEntityHook2 struct{}

func (*nullEntityHook2) Name() string              { return "null-hook2" }
func (*nullEntityHook2) Priority() int             { return 40 }
func (*nullEntityHook2) Run(_, _ *pb.Entity) error { return nil }

func goodEntityConstructor2(_ RefContext) (EntityHook, error) {
	return &nullEntityHook2{}, nil
}

func badEntityConstructor(_ RefContext) (EntityHook, error) {
	return nil, errors.New("initialization error")
}

type errorEntityHook struct{}

func (*errorEntityHook) Name() string              { return "error-hook" }
func (*errorEntityHook) Priority() int             { return 50 }
func (*errorEntityHook) Run(_, _ *pb.Entity) error { return errors.New("an Error") }
