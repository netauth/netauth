package util

import (
	pb "github.com/NetAuth/Protocol"
)

// LoadEntityBatch uses the specified loader to load a batch of
// entities as specified by the provided slice of IDs.  any failure in
// the loading of an entity will abort the entire load.
func LoadEntityBatch(ids []string, l loadEntityFunc) ([]*pb.Entity, error) {
	eSlice := []*pb.Entity{}

	for i := range ids {
		e, err := l(ids[i])
		if err != nil {
			return nil, err
		}
		eSlice = append(eSlice, e)
	}
	return eSlice, nil
}

// LoadGroupBatch functions identically to LoadEntityBatch, with
// appropriate substitutions made for groups.
func LoadGroupBatch(ids []string, l loadGroupFunc) ([]*pb.Group, error) {
	gSlice := []*pb.Group{}

	for i := range ids {
		g, err := l(ids[i])
		if err != nil {
			return nil, err
		}
		gSlice = append(gSlice, g)
	}
	return gSlice, nil
}
