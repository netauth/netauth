package tree

import (
	"sort"

	"github.com/NetAuth/NetAuth/internal/tree/errors"

	pb "github.com/NetAuth/Protocol"
)

// Run handles entity processor pipelines
func (ep *EntityProcessor) Run() (*pb.Entity, error) {
	for _, h := range ep.hooks {
		if err := h.Run(ep.Entity, ep.RequestData); err != nil {
			return nil, err
		}
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
