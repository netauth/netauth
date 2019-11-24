package util

import (
	"errors"

	"github.com/golang/protobuf/proto"

	pb "github.com/netauth/protocol"
)

func evilEntityLoader(_ string) (*pb.Entity, error) {
	return nil, errors.New("some obscure loading error")
}

func evilEntityIDFinder() ([]string, error) {
	return nil, errors.New("some obscure iterator error")
}

func okEntityLoader(s string) (*pb.Entity, error) {
	switch s {
	case "entity1":
		return &pb.Entity{
			ID:     proto.String("entity1"),
			Number: proto.Int32(2),
		}, nil
	case "entity2":
		return &pb.Entity{
			ID:     proto.String("entity2"),
			Number: proto.Int32(1),
		}, nil
	default:
		return nil, errors.New("wat")
	}
}

func okEntityIDFinder() ([]string, error) {
	return []string{"entity1", "entity2"}, nil
}

func evilGroupLoader(_ string) (*pb.Group, error) {
	return nil, errors.New("some obscure loading error")
}

func evilGroupIDFinder() ([]string, error) {
	return nil, errors.New("some obscure iterator error")
}

func okGroupLoader(s string) (*pb.Group, error) {
	switch s {
	case "group1":
		return &pb.Group{
			Name:   proto.String("group1"),
			Number: proto.Int32(2),
		}, nil
	case "group2":
		return &pb.Group{
			Name:   proto.String("group2"),
			Number: proto.Int32(1),
		}, nil
	default:
		return nil, errors.New("wat")
	}
}

func okGroupIDFinder() ([]string, error) {
	return []string{"group1", "group2"}, nil
}
