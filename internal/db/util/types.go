package util

import (
	pb "github.com/NetAuth/Protocol"
)

type loadEntityFunc func(string) (*pb.Entity, error)
type entityIDsFunc func() ([]string, error)

type loadGroupFunc func(string) (*pb.Group, error)
type groupNamesFunc func() ([]string, error)
