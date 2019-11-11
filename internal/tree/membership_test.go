package tree

import (
	"testing"

	"github.com/golang/protobuf/proto"

	pb "github.com/NetAuth/Protocol"
)

func TestDedupEntityList(t *testing.T) {
	eList := []*pb.Entity{
		{ID: proto.String("entity1")},
		{ID: proto.String("entity2")},
		{ID: proto.String("entity1")},
		{ID: proto.String("entity3")},
		{ID: proto.String("entity4")},
	}

	list := dedupEntityList(eList)

	if len(list) != 4 {
		t.Error("List size is wrong after deduping")
	}
}

func TestEntityListDifference(t *testing.T) {
	list1 := []*pb.Entity{
		{ID: proto.String("entity1")},
		{ID: proto.String("entity2")},
		{ID: proto.String("entity3")},
		{ID: proto.String("entity4")},
	}

	list2 := []*pb.Entity{
		{ID: proto.String("entity1")},
		{ID: proto.String("entity2")},
		{ID: proto.String("entity4")},
	}

	result := entityListDifference(list1, list2)

	if len(result) != 1 || result[0].GetID() != "entity3" {
		t.Error("List difference incorrect")
	}
}
