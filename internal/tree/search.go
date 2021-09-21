package tree

import (
	"context"

	"github.com/netauth/netauth/internal/db"

	pb "github.com/netauth/protocol"
)

// SearchGroups returns a list of groups filtered by the search
// criteria.
func (m *Manager) SearchGroups(ctx context.Context, r db.SearchRequest) ([]*pb.Group, error) {
	return m.db.SearchGroups(ctx, r)
}

// SearchEntities returns a list of entities filtered by the search
// criteria.
func (m *Manager) SearchEntities(ctx context.Context, r db.SearchRequest) ([]*pb.Entity, error) {
	entities, err := m.db.SearchEntities(ctx, r)
	if err != nil {
		return nil, err
	}

	out := make([]*pb.Entity, len(entities))
	for i := range entities {
		out[i] = safeCopyEntity(entities[i])
	}
	return out, nil
}
