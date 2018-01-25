package entity_manager

import (
	"github.com/golang/protobuf/proto"

	pb "github.com/NetAuth/NetAuth/pkg/proto"
)

// safeCopyEntity makes a copy of the entity provided but removes
// fields that are related to security.  This permits the entity that
// is returned to be handed off outside the server.
func safeCopyEntity(e *pb.Entity) (*pb.Entity, error) {
	// Marshal the proto to get a pure data representation of it.
	data, err := proto.Marshal(e)
	if err != nil {
		return nil, err
	}

	// Unmarshaling here ensures that the new entity has no
	// connection to the old one.
	ne := &pb.Entity{}
	if err := proto.Unmarshal(data, ne); err != nil {
		return nil, err
	}

	// Before returning, fields related to security are nulled out
	// so that they aren't available in the returned copy.  At
	// least not available in a meaningful sense.
	ne.Secret = proto.String("<REDACTED>")

	return ne, nil
}
