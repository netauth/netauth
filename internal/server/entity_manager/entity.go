package entity_manager

import (
	"log"

	"github.com/golang/protobuf/proto"
	"golang.org/x/crypto/bcrypt"

	"github.com/NetAuth/NetAuth/pkg/errors"
	pb "github.com/NetAuth/NetAuth/proto"
)

// New returns an initialized EMDataStore on to which all other
// functions are bound.
func New(db *EMDiskInterface) *EMDataStore {
	x := EMDataStore{}
	x.eByID = make(map[string]*pb.Entity)
	x.eByUIDNumber = make(map[int32]*pb.Entity)
	x.bootstrap_done = false
	if db != nil {
		x.db = *db
	}
	log.Println("Initialized new Entity Manager")

	if x.db == nil {
		log.Println("Entity Manager persistence layer is not available!")
	}

	return &x
}

// nextUIDNumber computes the next available uidNumber to be assigned.
// This allows a NewEntity request to be made with the uidNumber field
// unset.
func (emds *EMDataStore) nextUIDNumber() int32 {
	var largest int32 = 0

	// Iterate over the entities and return the largest ID found
	// +1.  This allows them to be in any order or have IDs
	// missing in the middle and still work.  Though an
	// inefficient search this is worst case O(N) and happens in
	// memory anyway.
	for i := range emds.eByID {
		if emds.eByID[i].GetUidNumber() > largest {
			largest = emds.eByID[i].GetUidNumber()
		}
	}

	return largest + 1
}

// newEntity creates a new entity given an ID, uidNumber, and secret.
// Its not necessary to set the secret upon creation and it can be set
// later.  If not set on creaction then the entity will not be usable.
// uidNumber must be a unique positive integer.  Because these are
// generally allocated in sequence the special value '-1' may be
// specified which will select the next available number.
func (emds *EMDataStore) newEntity(ID string, uidNumber int32, secret string) error {
	// Does this entity exist already?
	_, exists := emds.eByID[ID]
	if exists {
		log.Printf("Entity with ID '%s' already exists!", ID)
		return errors.E_DUPLICATE_ID
	}
	_, exists = emds.eByUIDNumber[uidNumber]
	if exists {
		log.Printf("Entity with uidNumber '%d' already exists!", uidNumber)
		return errors.E_DUPLICATE_UIDNUMBER
	}

	// Were we given a specific uidNumber?
	if uidNumber == -1 {
		// -1 is a sentinel value that tells us to pick the
		// next available number and assign it.
		uidNumber = emds.nextUIDNumber()
	}

	// Ok, they don't exist so we'll make them exist now
	newEntity := &pb.Entity{
		ID:        &ID,
		UidNumber: &uidNumber,
		Secret:    &secret,
		Meta:      &pb.EntityMeta{},
	}

	// Add this entity to the in-memory listings
	emds.eByID[ID] = newEntity
	emds.eByUIDNumber[uidNumber] = newEntity

	// Save the entity if the persistence layer is available.
	if emds.db != nil {
		if err := emds.db.SaveEntity(newEntity); err != nil {
			return err
		}
	}

	// Now we set the entity secret, this could be inlined, but
	// having it in the seperate func (emds *EMDataStore)tion makes resetting the
	// secret trivial.
	if err := emds.setEntitySecretByID(ID, secret); err != nil {
		return err
	}

	// Successfully created we now return no errors
	log.Printf("Created entity '%s'", ID)

	return nil
}

// NewEntity is a public function which adds a new entity on behalf of
// another one.  The requesting entity must be able to validate its
// identity and posses the appropriate capability to add a new entity
// to the system.
func (emds *EMDataStore) NewEntity(requestID, requestSecret, newID string, newUIDNumber int32, newSecret string) error {
	// Validate that the entity is real and permitted to perform
	// this action.
	if err := emds.validateEntityCapabilityAndSecret(requestID, requestSecret, "CREATE_ENTITY"); err != nil {
		return err
	}

	// The entity is who they say they are and has the appropriate
	// capability, time to actually create the new entity.
	if err := emds.newEntity(newID, newUIDNumber, newSecret); err != nil {
		return err
	}
	return nil
}

// NewBootstrapEntity is a function that can be called during the
// startup of the srever to create an entity that has the appropriate
// authority to create more entities and otherwise manage the server.
// This can only be called once during startup, attepts to call it
// again will result in no change.  The bootstrap user will always get
// the next available number which in most cases will be 1.
func (emds *EMDataStore) MakeBootstrap(ID string, secret string) {
	if emds.bootstrap_done {
		return
	}

	// In some cases if there is an existing system that has no
	// admin, it is necessary to confer bootstrap powers to an
	// existing user.  In that case they are just selected and
	// then provided the GLOBAL_ROOT capability.
	e, err := emds.getEntityByID(ID)
	if err != nil {
		log.Printf("No entity with ID '%s' exists!  Creating...", ID)
	}

	// This is not a normal Go way of doing this, but this
	// function has two possible success cases, the flow may jump
	// in here and return if there is an existing entity to get
	// root powers.
	if e != nil {
		emds.setEntityCapability(e, "GLOBAL_ROOT")
		emds.bootstrap_done = true
		return
	}

	// Even in the bootstrap case its still possible this can
	// fail, in that case its useful to have the error.
	if err := emds.newEntity(ID, -1, secret); err != nil {
		log.Printf("Could not create bootstrap user! (%s)", err)
	}
	if err := emds.setEntityCapabilityByID(ID, "GLOBAL_ROOT"); err != nil {
		log.Printf("Couldn't provide root authority! (%s)", err)
	}

	emds.bootstrap_done = true
}

// DisableBootstrap disables the ability to bootstrap after the
// opportunity to do so has passed.
func (emds *EMDataStore) DisableBootstrap() {
	emds.bootstrap_done = true
}

// loadFromDisk attempts to load an entity from the disk.  This can be
// useful if another server wrote it to disk but its not in the cache
// yet.
func (emds *EMDataStore) loadFromDisk(ID string) error {
	// If the persistence layer isn't available, just return.
	// This isn't necessarily an error, but no data was available.
	if emds.db == nil {
		return nil
	}

	e, err := emds.db.LoadEntity(ID)
	if err != nil {
		return err
	}

	emds.eByID[e.GetID()] = e
	emds.eByUIDNumber[e.GetUidNumber()] = e

	log.Printf("Loaded '%s' from disk", e.GetID())
	return nil
}

// getEntityByID returns a pointer to an Entity struct and an error
// value.  The error value will either be errors.E_NO_ENTITY if the
// requested value did not match, or will be nil where an entity is
// returned.  The string must be a complete match for the entity name
// being requested.
func (emds *EMDataStore) getEntityByID(ID string) (*pb.Entity, error) {
	e, ok := emds.eByID[ID]
	if !ok {
		// Attempt to load the entity if the persistence layer
		// is available.
		if emds.db != nil {
			if err := emds.loadFromDisk(ID); err != nil {
				return nil, err
			}
		}
	}

	// Try again after potentially having loaded the entity
	e, ok = emds.eByID[ID]
	if !ok {
		return nil, errors.E_NO_ENTITY
	}
	return e, nil
}

// DeleteEntityByID deletes the named entity.  This func (emds
// *EMDataStore)tion will delete the entity in a non-atomic way, but
// will ensure that the entity cannot be authenticated with before
// returning.  If the named ID does not exist the function will return
// errors.E_NO_ENTITY, in all other cases nil is returned.
func (emds *EMDataStore) deleteEntityByID(ID string) error {
	// Drop the entity from the in memory storage.
	if err := emds.dropEntity(ID); err != nil {
		return err
	}

	// Delete the entity if the persistence layer is available.
	if emds.db != nil {
		if err := emds.db.DeleteEntity(ID); err != nil {
			return err
		}
	}

	log.Printf("Deleted entity '%s'", ID)

	return nil
}

func (emds *EMDataStore) DeleteEntity(requestID string, requestSecret string, deleteID string) error {
	// Validate that the entity is real and permitted to perform
	// this action.
	if err := emds.validateEntityCapabilityAndSecret(requestID, requestSecret, "DELETE_ENTITY"); err != nil {
		return err
	}

	// Delete the requested entity
	return emds.deleteEntityByID(deleteID)
}

// dropEntity is similar to delete, but just removes the entity from
// the in memory data store.
func (emds *EMDataStore) dropEntity(ID string) error {
	e, err := emds.getEntityByID(ID)
	if err != nil {
		return err
	}

	delete(emds.eByID, e.GetID())
	delete(emds.eByUIDNumber, e.GetUidNumber())

	return nil
}

// checkCapability is a helper function which allows a method to
// quickly check for a capability on an entity.  This check only looks
// for capabilities that an entity has directly, not any which may be
// conferred to it by group membership.
func (emds *EMDataStore) checkEntityCapability(e *pb.Entity, c string) error {
	for _, a := range e.Meta.Capabilities {
		if a == pb.Capability_GLOBAL_ROOT {
			return nil
		}

		if a == pb.Capability(pb.Capability_value[c]) {
			return nil
		}
	}
	return errors.E_ENTITY_UNQUALIFIED
}

// checkCapabilityByID is a convenience func (emds *EMDataStore)tion which performs the
// query to retrieve the entity itself, rather than requirin the
// caller to produce the pointer to the entity.
func (emds *EMDataStore) checkEntityCapabilityByID(ID string, c string) error {
	e, err := emds.getEntityByID(ID)
	if err != nil {
		return err
	}

	return emds.checkEntityCapability(e, c)
}

// SetCapability sets a capability on an entity.  The set operation is
// idempotent.
func (emds *EMDataStore) setEntityCapability(e *pb.Entity, c string) error {
	// If no capability was supplied, bail out.
	if len(c) == 0 {
		return nil
	}

	cap := pb.Capability(pb.Capability_value[c])

	for _, a := range e.Meta.Capabilities {
		if a == cap {
			// The entity already has this capability
			// directly, don't add it again.
			return nil
		}
	}

	e.Meta.Capabilities = append(e.Meta.Capabilities, cap)

	// Save the entity if the persistence layer is available.
	if emds.db != nil {
		if err := emds.db.SaveEntity(e); err != nil {
			return err
		}
	}

	log.Printf("Set capability %s on entity '%s'", c, e.GetID())
	return nil
}

// SetEntityCapabilityByID is a convenience function to get the entity
// and hand it off to the actual setEntityCapability function
func (emds *EMDataStore) setEntityCapabilityByID(ID string, c string) error {
	e, err := emds.getEntityByID(ID)
	if err != nil {
		return err
	}

	return emds.setEntityCapability(e, c)
}

// SetEntitySecretByID sets the secret on a given entity using the
// bcrypt secure hashing algorithm.
func (emds *EMDataStore) setEntitySecretByID(ID string, secret string) error {
	e, err := emds.getEntityByID(ID)
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

	// Save the entity if the persistence layer is available.
	if emds.db != nil {
		if err := emds.db.SaveEntity(e); err != nil {
			return err
		}
	}

	log.Printf("Secret set for '%s'", e.GetID())
	return nil
}

// ChangeSecret is a publicly available function to change an entity
// secret.  This function requires either the CHANGE_ENTITY_SECRET
// capability or the entity to be requesting the change for itself.
func (emds *EMDataStore) ChangeSecret(ID string, secret string, changeID string, changeSecret string) error {
	// If the entity isn't the one requesting the change then
	// extra capabilities are required.
	if ID != changeID {
		if err := emds.validateEntityCapabilityAndSecret(ID, secret, "CHANGE_ENTITY_SECRET"); err != nil {
			return err
		}
	} else {
		if err := emds.ValidateSecret(ID, secret); err != nil {
			return err
		}
	}

	// At this point the entity is either the one that we're
	// changing the secret for or is the one that is allowed to
	// change the secrets of others.
	if err := emds.setEntitySecretByID(changeID, changeSecret); err != nil {
		return err
	}

	// At this point the secret has been changed.
	return nil
}

// ValidateSecret validates the identity of an entity by
// validating the authenticating entity with the secret.
func (emds *EMDataStore) ValidateSecret(ID string, secret string) error {
	e, err := emds.getEntityByID(ID)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(*e.Secret), []byte(secret))
	if err != nil {
		// This is strictly not in the style of go, but this
		// is the best place to put this log message so that
		// it works like all the others.
		log.Printf("Failed to authenticate '%s'", e.GetID())
		return errors.E_ENTITY_BADAUTH
	}
	log.Printf("Successfully authenticated '%s'", e.GetID())

	return nil
}

// validateEntityCapabilityAndSecret validates an entitity is who they
// say they are and that they have a named capability.  This is a
// convenience function and simply calls and aggregates responses from
// other functions which perform the actual checks.
func (emds *EMDataStore) validateEntityCapabilityAndSecret(ID string, secret string, capability string) error {
	// First validate the entity identity.
	if err := emds.ValidateSecret(ID, secret); err != nil {
		return err
	}

	// Then validate the entity capability.
	if err := emds.checkEntityCapabilityByID(ID, capability); err != nil {
		return err
	}

	// todo(maldridge) When groups have capabilities this may be
	// checked here as well.

	// Entity is who they say they are and has the specified capability.
	return nil
}

// GetEntity returns an entity to the caller after first making a safe
// copy of it to remove secure fields.
func (emds *EMDataStore) GetEntity(ID string) (*pb.Entity, error) {
	// e will be the direct internal copy, we can't give this back
	// though since it has secrets embedded.
	e, err := emds.getEntityByID(ID)
	if err != nil {
		return nil, err
	}

	// The safeCopyEntity will return the entity without secrets
	// in it, as well as an error if there were problems
	// marshaling the proto back and forth.
	return safeCopyEntity(e)
}

func (emds *EMDataStore) updateEntityMeta(e *pb.Entity, newMeta *pb.EntityMeta) error {
	// get the existing metadata
	meta := e.GetMeta()

	// some fields must not be merged in, so we make sure that
	// they're nulled out here
	newMeta.Capabilities = nil
	newMeta.Groups = nil

	// now we can merge the changes, this happens on the live tree
	// and doesn't require recomputing anything since its a change
	// at the leaves since the groups are not permitted to change
	// by this API.
	proto.Merge(meta, newMeta)

	// Save the entity if the persistence layer is available.
	if emds.db != nil {
		if err := emds.db.SaveEntity(e); err != nil {
			return err
		}
	}

	log.Printf("Updated metadata for '%s'", e.GetID())
	return nil
}

func (emds *EMDataStore) UpdateEntityMeta(requestID, requestSecret, modEntityID string, newMeta *pb.EntityMeta) error {
	// An entity can change its own metadata, but modifying other
	// entities requires an administrative capability.
	if modEntityID != requestID {
		if err := emds.validateEntityCapabilityAndSecret(requestID, requestSecret, "MODIFY_ENTITY_META"); err != nil {
			log.Printf("UpdateEntityMeta denied to '%s'", requestID)
			return err
		}
	} else {
		if err := emds.ValidateSecret(requestID, requestSecret); err != nil {
			log.Printf("UpdateEntityMeta denied to '%s'", requestID)
			return err
		}
	}

	// Get the actual entity that the metadata is being updated
	// on.
	e, err := emds.getEntityByID(modEntityID)
	if err != nil {
		return err
	}

	// Run the update
	return emds.updateEntityMeta(e, newMeta)
}
