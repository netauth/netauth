package tree

import (
	"errors"
	"testing"

	pb "github.com/NetAuth/Protocol"
)

func resetGConstructorMap() {
	gHookConstructors = make(map[string]GroupHookConstructor)
}

func TestGPRegisterAndInitialize(t *testing.T) {
	resetGConstructorMap()
	defer resetGConstructorMap()

	RegisterGroupHookConstructor("null-hook", goodGroupConstructor)
	RegisterGroupHookConstructor("null-hook", goodGroupConstructor)

	if len(gHookConstructors) != 1 {
		t.Error("Duplicate hook registered")
	}

	RegisterGroupHookConstructor("bad-hook", badGroupConstructor)

	if len(gHookConstructors) != 2 {
		t.Error("bad-hook wasn't registered")
	}

	em := Manager{
		groupProcessorHooks: make(map[string]GroupProcessorHook),
	}

	em.InitializeGroupHooks()
	if len(em.groupProcessorHooks) != 1 {
		t.Error("bad-hook was initialized")
	}
}

func TestGPInitializeChainsOK(t *testing.T) {
	resetGConstructorMap()
	defer resetGConstructorMap()

	RegisterGroupHookConstructor("null-hook", goodGroupConstructor)
	RegisterGroupHookConstructor("null-hook2", goodGroupConstructor2)
	em := Manager{
		groupProcessorHooks: make(map[string]GroupProcessorHook),
		groupProcesses:      make(map[string][]GroupProcessorHook),
	}
	em.InitializeGroupHooks()

	c := map[string][]string{
		"TEST": []string{"null-hook", "null-hook2"},
	}

	if err := em.InitializeGroupChains(c); err != nil {
		t.Error(err)
	}
}

func TestGPInitializeBadHook(t *testing.T) {
	resetGConstructorMap()
	defer resetGConstructorMap()

	em := Manager{
		groupProcessorHooks: make(map[string]GroupProcessorHook),
		groupProcesses:      make(map[string][]GroupProcessorHook),
	}
	em.InitializeGroupHooks()

	c := map[string][]string{
		"TEST": []string{"unknown-hook"},
	}

	if err := em.InitializeGroupChains(c); err != ErrUnknownHook {
		t.Error(err)
	}
}

func TestGPFetchHooksOK(t *testing.T) {
	hookmap := map[string][]GroupProcessorHook{
		"TEST": []GroupProcessorHook{
			&nullGroupHook{},
		},
	}

	ep := GroupProcessor{}

	if err := ep.FetchHooks("TEST", hookmap); err != nil {
		t.Error(err)
	}
}

func TestGPFetchHooksBadChain(t *testing.T) {
	hookmap := map[string][]GroupProcessorHook{}
	ep := GroupProcessor{}

	if err := ep.FetchHooks("MISSING", hookmap); err != ErrUnknownHookChain {
		t.Error(err)
	}
}

func TestGPFetchHooksEmptyChain(t *testing.T) {
	hookmap := map[string][]GroupProcessorHook{
		"EMPTY": []GroupProcessorHook{},
	}
	ep := GroupProcessor{}

	if err := ep.FetchHooks("EMPTY", hookmap); err != ErrEmptyHookChain {
		t.Error(err)
	}
}

func TestGPRunOK(t *testing.T) {
	ep := GroupProcessor{
		hooks: []GroupProcessorHook{&nullGroupHook{}},
	}

	if _, err := ep.Run(); err != nil {
		t.Error(err)
	}
}

func TestGPRunHookError(t *testing.T) {
	ep := GroupProcessor{
		hooks: []GroupProcessorHook{&errorGroupHook{}},
	}

	if _, err := ep.Run(); err == nil {
		t.Error("Hook runtime ate an error...")
	}
}

type nullGroupHook struct{}

func (*nullGroupHook) Name() string             { return "null-hook" }
func (*nullGroupHook) Priority() int            { return 50 }
func (*nullGroupHook) Run(_, _ *pb.Group) error { return nil }

func goodGroupConstructor(_ RefContext) (GroupProcessorHook, error) {
	return &nullGroupHook{}, nil
}

type nullGroupHook2 struct{}

func (*nullGroupHook2) Name() string             { return "null-hook2" }
func (*nullGroupHook2) Priority() int            { return 40 }
func (*nullGroupHook2) Run(_, _ *pb.Group) error { return nil }

func goodGroupConstructor2(_ RefContext) (GroupProcessorHook, error) {
	return &nullGroupHook2{}, nil
}

func badGroupConstructor(_ RefContext) (GroupProcessorHook, error) {
	return nil, errors.New("initialization error")
}

type errorGroupHook struct{}

func (*errorGroupHook) Name() string             { return "error-hook" }
func (*errorGroupHook) Priority() int            { return 50 }
func (*errorGroupHook) Run(_, _ *pb.Group) error { return errors.New("an Error") }
