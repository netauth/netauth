package tree

import (
	"errors"
	"testing"

	pb "github.com/NetAuth/Protocol"
)

func TestEPRegisterAndInitialize(t *testing.T) {
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
		entityProcessorHooks: make(map[string]EntityProcessorHook),
	}

	em.InitializeEntityHooks()
	if len(em.entityProcessorHooks) != 1 {
		t.Error("bad-hook was initialized")
	}
}

func TestEPInitializeChainsOK(t *testing.T) {
	RegisterEntityHookConstructor("null-hook", goodEntityConstructor)
	RegisterEntityHookConstructor("null-hook2", goodEntityConstructor2)
	em := Manager{
		entityProcessorHooks: make(map[string]EntityProcessorHook),
		entityProcesses:      make(map[string][]EntityProcessorHook),
	}
	em.InitializeEntityHooks()

	c := map[string][]string{
		"TEST": []string{"null-hook", "null-hook2"},
	}

	if err := em.InitializeEntityChains(c); err != nil {
		t.Error(err)
	}
}

func TestEPInitializeBadHook(t *testing.T) {
	em := Manager{
		entityProcessorHooks: make(map[string]EntityProcessorHook),
		entityProcesses:      make(map[string][]EntityProcessorHook),
	}
	em.InitializeEntityHooks()

	c := map[string][]string{
		"TEST": []string{"unknown-hook"},
	}

	if err := em.InitializeEntityChains(c); err != ErrUnknownHook {
		t.Error(err)
	}
}

func TestEPFetchHooksOK(t *testing.T) {
	hookmap := map[string][]EntityProcessorHook{
		"TEST": []EntityProcessorHook{
			&nullEntityHook{},
		},
	}

	ep := EntityProcessor{}

	if err := ep.FetchHooks("TEST", hookmap); err != nil {
		t.Error(err)
	}
}

func TestEPFetchHooksBadChain(t *testing.T) {
	hookmap := map[string][]EntityProcessorHook{}
	ep := EntityProcessor{}

	if err := ep.FetchHooks("MISSING", hookmap); err != ErrUnknownHookChain {
		t.Error(err)
	}
}

func TestEPFetchHooksEmptyChain(t *testing.T) {
	hookmap := map[string][]EntityProcessorHook{
		"EMPTY": []EntityProcessorHook{},
	}
	ep := EntityProcessor{}

	if err := ep.FetchHooks("EMPTY", hookmap); err != ErrEmptyHookChain {
		t.Error(err)
	}
}

func TestEPRunOK(t *testing.T) {
	ep := EntityProcessor{
		hooks: []EntityProcessorHook{&nullEntityHook{}},
	}

	if _, err := ep.Run(); err != nil {
		t.Error(err)
	}
}

func TestEPRunHookError(t *testing.T) {
	ep := EntityProcessor{
		hooks: []EntityProcessorHook{&errorEntityHook{}},
	}

	if _, err := ep.Run(); err == nil {
		t.Error("Hook runtime ate an error...")
	}
}

type nullEntityHook struct{}

func (*nullEntityHook) Name() string              { return "null-hook" }
func (*nullEntityHook) Priority() int             { return 50 }
func (*nullEntityHook) Run(_, _ *pb.Entity) error { return nil }

func goodEntityConstructor(_ RefContext) (EntityProcessorHook, error) {
	return &nullEntityHook{}, nil
}

type nullEntityHook2 struct{}

func (*nullEntityHook2) Name() string              { return "null-hook2" }
func (*nullEntityHook2) Priority() int             { return 40 }
func (*nullEntityHook2) Run(_, _ *pb.Entity) error { return nil }

func goodEntityConstructor2(_ RefContext) (EntityProcessorHook, error) {
	return &nullEntityHook2{}, nil
}

func badEntityConstructor(_ RefContext) (EntityProcessorHook, error) {
	return nil, errors.New("initialization error")
}

type errorEntityHook struct{}

func (*errorEntityHook) Name() string              { return "error-hook" }
func (*errorEntityHook) Priority() int             { return 50 }
func (*errorEntityHook) Run(_, _ *pb.Entity) error { return errors.New("an Error") }
