package util

import (
	pb "github.com/NetAuth/Protocol"
)

type loadEntityFunc func(string) (*pb.Entity, error)
type entityIDsFunc func() ([]string, error)

// NextEntityNumber computes and returns the next unnassigned number
// in the entity space.
func NextEntityNumber(l loadEntityFunc, ids entityIDsFunc) (int32, error) {
	var largest int32

	// Iterate over the entities and return the largest ID found
	// +1.  This allows them to be in any order or have IDs
	// missing in the middle and still work.  Though an
	// inefficient search this is worst case O(N) and happens only
	// on provisioning a new entry in the database.
	el, err := ids()
	if err != nil {
		return 0, err
	}

	for _, en := range el {
		e, err := l(en)
		if err != nil {
			return 0, err
		}
		if e.GetNumber() > largest {
			largest = e.GetNumber()
		}
	}

	return largest + 1, nil
}

type loadGroupFunc func(string) (*pb.Group, error)
type groupNamesFunc func() ([]string, error)

// NextGroupNumber computes the enxt available group number and
// returns it.
func NextGroupNumber(lf loadGroupFunc, names groupNamesFunc) (int32, error) {
	var largest int32

	l, err := names()
	if err != nil {
		return 0, err
	}
	for _, i := range l {
		g, err := lf(i)
		if err != nil {
			return 0, err
		}
		if g.GetNumber() > largest {
			largest = g.GetNumber()
		}
	}

	return largest + 1, nil
}
