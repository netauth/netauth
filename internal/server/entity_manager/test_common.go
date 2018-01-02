package entity_manager

import (
	pb "github.com/NetAuth/NetAuth/proto"
)

func resetMap() {
	eByID = make(map[string]*pb.Entity)
	eByUIDNumber = make(map[int32]*pb.Entity)
}
