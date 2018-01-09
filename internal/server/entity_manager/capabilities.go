package entity_manager

import (
	"log"

	"golang.org/x/crypto/bcrypt"

	pb "github.com/NetAuth/NetAuth/proto"
)

// checkCapability is a helper function which allows a method to
// quickly check for a capability on an entity.  This check only looks
// for capabilities that an entity has directly, not any which may be
// conferred to it by group membership.
func checkEntityCapability(e *pb.Entity, c string) error {
	for _, a := range e.Meta.Capabilities {
		if a == pb.Capability_GLOBAL_ROOT {
			return nil
		}

		if a == pb.Capability(pb.Capability_value[c]) {
			return nil
		}
	}
	return E_ENTITY_UNQUALIFIED
}

// checkCapabilityByID is a convenience function which performs the
// query to retrieve the entity itself, rather than requirin the
// caller to produce the pointer to the entity.
func checkEntityCapabilityByID(ID string, c string) error {
	e, err := getEntityByID(ID)
	if err != nil {
		return err
	}

	return checkEntityCapability(e, c)
}

// SetCapability sets a capability on an entity.  The set operation is
// idempotent.
func setEntityCapability(e *pb.Entity, c string) {
	// If no capability was supplied, bail out.
	if len(c) == 0 {
		return
	}

	cap := pb.Capability(pb.Capability_value[c])

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

// SetEntityCapabilityByID is a convenience function to get the entity
// and hand it off to the actual setEntityCapability function
func setEntityCapabilityByID(ID string, c string) error {
	e, err := getEntityByID(ID)
	if err != nil {
		return err
	}

	setEntityCapability(e, c)
	return nil
}

// SetEntitySecretByID sets the secret on a given entity using the
// bcrypt secure hashing algorithm.
func setEntitySecretByID(ID string, secret string) error {
	e, err := getEntityByID(ID)
	if err != nil {
		return err
	}

	// todo(maldridge) This is currently configured to use the
	// minimum cost hash.  THIS IS NOT SECURE.  Before releasing
	// 1.0 this must be changes to pull this value either
	// dynamically or from a config file somewhere.
	hash, err := bcrypt.GenerateFromPassword([]byte(secret), 0)
	if err != nil {
		return err
	}
	hashedSecret := string(hash[:])
	e.Secret = &hashedSecret

	log.Printf("Secret set for '%s'", e.GetID())
	return nil
}

// ChangeSecret is a publicly available function to change an entity
// secret.  This function requires either the CHANGE_ENTITY_SECRET
// capability or the entity to be requesting the change for itself.
func ChangeSecret(ID string, secret string, changeID string, changeSecret string) error {
	// If the entity isn't the one requesting the change then
	// extra capabilities are required.
	if ID != changeID {
		if err := validateEntityCapabilityAndSecret(ID, secret, "CHANGE_ENTITY_SECRET"); err != nil {
			return err
		}
	} else {
		if err := ValidateEntitySecretByID(ID, secret); err != nil {
			return err
		}
	}

	// At this point the entity is either the one that we're
	// changing the secret for or is the one that is allowed to
	// change the secrets of others.
	if err := setEntitySecretByID(changeID, changeSecret); err != nil {
		return err
	}

	// At this point the secret has been changed.
	return nil
}

// ValidateEntitySecretByID validates the identity of an entity by
// validating the authenticating entity with the secret.
func ValidateEntitySecretByID(ID string, secret string) error {
	e, err := getEntityByID(ID)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(*e.Secret), []byte(secret))
	if err != nil {
		// This is strictly not in the style of go, but this
		// is the best place to put this log message so that
		// it works like all the others.
		log.Printf("Failed to authenticate '%s'", e.GetID())
		return E_ENTITY_BADAUTH
	}
	log.Printf("Successfully authenticated '%s'", e.GetID())

	return nil
}

// validateEntityCapabilityAndSecret validates an entitity is who they
// say they are and that they have a named capability.  This is a
// convenience function and simply calls and aggregates responses from
// other functions which perform the actual checks.
func validateEntityCapabilityAndSecret(ID string, secret string, capability string) error {
	// First validate the entity identity.
	if err := ValidateEntitySecretByID(ID, secret); err != nil {
		return err
	}

	// Then validate the entity capability.
	if err := checkEntityCapabilityByID(ID, capability); err != nil {
		return err
	}

	// todo(maldridge) When groups have capabilities this may be
	// checked here as well.

	// Entity is who they say they are and has the specified capability.
	return nil
}
