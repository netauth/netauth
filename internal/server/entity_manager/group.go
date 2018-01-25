package entity_manager

import (
	"github.com/NetAuth/NetAuth/pkg/errors"
	pb "github.com/NetAuth/NetAuth/proto"
)

// listMembers takes a group ID in and returns a slice of entities
// that are in that group.
func (emds *EMDataStore) listMembers(groupID string) ([]*pb.Entity, error) {
	// 'ALL' is a special groupID which returns everything, this
	// isn't a group that exists in a real sense, it just serves
	// to return a global list as a convenience.
	if groupID == "ALL" {
		var entities []*pb.Entity
		for _, e := range emds.eByID {
			entities = append(entities, e)
		}
		return entities, nil
	}

	// No group matched (likely because no other group mechanisms
	// are implemented).
	return nil, errors.E_NO_GROUP
}

// ListMembers fulfills the same function as the private version of
// this function, but with one crucial difference, it produces copies
// of the entities that have the secret redacted.
func (emds *EMDataStore) ListMembers(groupID string) ([]*pb.Entity, error) {
	// This set of members has secrets and can't be returned.
	members, err := emds.listMembers(groupID)
	if err != nil {
		return nil, err
	}

	var safeMembers []*pb.Entity
	for _, e := range members {
		ne, err := safeCopyEntity(e)
		if err != nil {
			return nil, err
		}
		safeMembers = append(safeMembers, ne)
	}

	return safeMembers, nil
}
