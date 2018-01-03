package entity_manager

import (
	"log"

	pb "github.com/NetAuth/NetAuth/proto"
)

// checkCapability is a helper function which allows a method to
// quickly check for a capability on an entity.  This check only looks
// for capabilities that an entity has directly, not any which may be
// conferred to it by group membership.
func checkCapability(e *pb.Entity, c string) error {
	for _, a := range e.Meta.Capabilities {
		if a == pb.EntityMeta_GLOBAL_ROOT {
			return nil
		}

		if a == pb.EntityMeta_Capability(pb.EntityMeta_Capability_value[c]) {
			return nil
		}
	}
	return E_ENTITY_UNQUALIFIED
}

// checkCapabilityByID is a convenience function which performs the
// query to retrieve the entity itself, rather than requirin the
// caller to produce the pointer to the entity.
func checkCapabilityByID(ID string, c string) error {
	e, err := getEntityByID(ID)
	if err != nil {
		return err
	}

	return checkCapability(e, c)
}

// SetCapability sets a capability on an entity.  The set operation is
// idempotent.
func setCapability(e *pb.Entity, c string) {
	cap := pb.EntityMeta_Capability(pb.EntityMeta_Capability_value[c])

	for _, a := range e.Meta.Capabilities {
		if a == cap {
			// The entity already has this capability
			// directly, don't add it again.
			return
		}
	}
	
	e.Meta.Capabilities = append(e.Meta.Capabilities, cap)
	log.Printf("Set capability %s on entity '%s'", c, e.GetID())
}
