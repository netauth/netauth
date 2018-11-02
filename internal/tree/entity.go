package tree

import (
	"fmt"
	"log"
	"strings"

	"github.com/golang/protobuf/proto"

	pb "github.com/NetAuth/Protocol"
)

// NewEntity creates a new entity given an ID, number, and secret.
// Its not necessary to set the secret upon creation and it can be set
// later.  If not set on creation then the entity will not be usable.
// number must be a unique positive integer.  Because these are
// generally allocated in sequence the special value '-1' may be
// specified which will select the next available number.
func (m *Manager) NewEntity(ID string, number int32, secret string) error {
	// Does this entity exist already?
	if _, err := m.db.LoadEntity(ID); err == nil {
		log.Printf("Entity with ID '%s' already exists!", ID)
		return ErrDuplicateEntityID
	}

	// Were we given a specific number?
	if number == -1 {
		var err error
		// -1 is a sentinel value that tells us to pick the
		// next available number and assign it.
		number, err = m.nextUIDNumber()
		if err != nil {
			return err
		}
	}

	// Ok, they don't exist so we'll make them exist now
	newEntity := &pb.Entity{
		ID:     &ID,
		Number: &number,
		Secret: &secret,
		Meta:   &pb.EntityMeta{},
	}

	// Save the entity
	if err := m.db.SaveEntity(newEntity); err != nil {
		return err
	}

	// Now we set the entity secret, this could be inlined, but
	// having it in the separate function makes resetting the
	// secret trivial.
	if err := m.SetEntitySecretByID(ID, secret); err != nil {
		return err
	}

	// Successfully created we now return no errors
	log.Printf("Created entity '%s'", ID)

	return nil
}

// nextUIDNumber computes the next available number to be assigned.
// This allows a NewEntity request to be made with the number field
// unset.
func (m *Manager) nextUIDNumber() (int32, error) {
	var largest int32

	// Iterate over the entities and return the largest ID found
	// +1.  This allows them to be in any order or have IDs
	// missing in the middle and still work.  Though an
	// inefficient search this is worst case O(N) and happens only
	// on provisioning a new entry in the database.
	el, err := m.db.DiscoverEntityIDs()
	if err != nil {
		return 0, err
	}

	for _, en := range el {
		e, err := m.db.LoadEntity(en)
		if err != nil {
			return 0, err
		}
		if e.GetNumber() > largest {
			largest = e.GetNumber()
		}
	}

	return largest + 1, nil
}

// MakeBootstrap is a function that can be called during the startup
// of the srever to create an entity that has the appropriate
// authority to create more entities and otherwise manage the server.
// This can only be called once during startup, attepts to call it
// again will result in no change.  The bootstrap user will always get
// the next available number which in most cases will be 1.
func (m *Manager) MakeBootstrap(ID string, secret string) {
	if m.bootstrapDone {
		return
	}

	// In some cases if there is an existing system that has no
	// admin, it is necessary to confer bootstrap powers to an
	// existing user.  In that case they are just selected and
	// then provided the GLOBAL_ROOT capability.
	e, err := m.db.LoadEntity(ID)
	if err != nil {
		log.Printf("No entity with ID '%s' exists!  Creating...", ID)
	}

	// We need to force unlock an entity here in the event that
	// the global bootstrap user somehow got locked.
	if e.GetMeta().GetLocked() {
		log.Println("Unlocking Entity")
		e.Meta.Locked = proto.Bool(false)
		if err := m.db.SaveEntity(e); err != nil {
			log.Printf("Could not unlock bootstrap entity: %s", err)
		}
	}

	// This is not a normal Go way of doing this, but this
	// function has two possible success cases, the flow may jump
	// in here and return if there is an existing entity to get
	// root powers.
	if e != nil {
		m.setEntityCapability(e, "GLOBAL_ROOT")
		m.bootstrapDone = true
		return
	}

	// Even in the bootstrap case its still possible this can
	// fail, in that case its useful to have the error.
	if err := m.NewEntity(ID, -1, secret); err != nil {
		log.Printf("Could not create bootstrap user! (%s)", err)
	}
	if err := m.SetEntityCapabilityByID(ID, "GLOBAL_ROOT"); err != nil {
		log.Printf("Couldn't provide root authority! (%s)", err)
	}

	m.bootstrapDone = true
}

// DisableBootstrap disables the ability to bootstrap after the
// opportunity to do so has passed.
func (m *Manager) DisableBootstrap() {
	log.Println("Disabling bootstrap")
	m.bootstrapDone = true
	log.Println("Bootstrap disabled")
}

// DeleteEntityByID deletes the named entity.  This function will
// delete the entity in a non-atomic way, but will ensure that the
// entity cannot be authenticated with before returning.  If the named
// ID does not exist the function will return errors.E_NO_ENTITY, in
// all other cases nil is returned.
func (m *Manager) DeleteEntityByID(ID string) error {
	if err := m.db.DeleteEntity(ID); err != nil {
		return err
	}
	log.Printf("Deleted entity '%s'", ID)

	return nil
}

// SetCapability sets a capability on an entity.  The set operation is
// idempotent.
func (m *Manager) setEntityCapability(e *pb.Entity, c string) error {
	// If no capability was supplied, bail out.
	if len(c) == 0 {
		return ErrUnknownCapability
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

	if err := m.db.SaveEntity(e); err != nil {
		return err
	}

	log.Printf("Set capability %s on entity '%s'", c, e.GetID())
	return nil
}

// removeCapability removes a capability on an entity
func (m *Manager) removeEntityCapability(e *pb.Entity, c string) error {
	// If no capability was supplied, bail out.
	if len(c) == 0 {
		return ErrUnknownCapability
	}

	cap := pb.Capability(pb.Capability_value[c])
	var ncaps []pb.Capability

	for _, a := range e.Meta.Capabilities {
		if a == cap {
			continue
		}
		ncaps = append(ncaps, a)
	}

	e.Meta.Capabilities = ncaps

	if err := m.db.SaveEntity(e); err != nil {
		return err
	}

	log.Printf("Removed capability %s on entity '%s'", c, e.GetID())
	return nil
}

// SetEntityCapabilityByID is a convenience function to get the entity
// and hand it off to the actual setEntityCapability function
func (m *Manager) SetEntityCapabilityByID(ID string, c string) error {
	e, err := m.db.LoadEntity(ID)
	if err != nil {
		return err
	}

	return m.setEntityCapability(e, c)
}

// RemoveEntityCapabilityByID is a convenience function to get the entity
// and hand it off to the actual removeEntityCapability function
func (m *Manager) RemoveEntityCapabilityByID(ID string, c string) error {
	e, err := m.db.LoadEntity(ID)
	if err != nil {
		return err
	}

	return m.removeEntityCapability(e, c)
}

// SetEntitySecretByID sets the secret on a given entity using the
// crypto interface.
func (m *Manager) SetEntitySecretByID(ID string, secret string) error {
	e, err := m.db.LoadEntity(ID)
	if err != nil {
		return err
	}

	ssecret, err := m.crypto.SecureSecret(secret)
	if err != nil {
		return err
	}
	e.Secret = &ssecret

	if err := m.db.SaveEntity(e); err != nil {
		return err
	}

	log.Printf("Secret set for '%s'", e.GetID())
	return nil
}

// ValidateSecret validates the identity of an entity by
// validating the authenticating entity with the secret.
func (m *Manager) ValidateSecret(ID string, secret string) error {
	e, err := m.db.LoadEntity(ID)
	if err != nil {
		return err
	}

	// Locked entities can't validate.
	if e.GetMeta().GetLocked() {
		return ErrEntityLocked
	}

	err = m.crypto.VerifySecret(secret, *e.Secret)
	if err != nil {
		log.Printf("Failed to authenticate '%s'", e.GetID())
		return err
	}
	log.Printf("Successfully authenticated '%s'", e.GetID())

	return nil
}

// GetEntity returns an entity to the caller after first making a safe
// copy of it to remove secure fields.
func (m *Manager) GetEntity(ID string) (*pb.Entity, error) {
	// e will be the direct internal copy, we can't give this back
	// though since it has secrets embedded.
	e, err := m.db.LoadEntity(ID)
	if err != nil {
		return nil, err
	}

	// The safeCopyEntity will return the entity without secrets
	// in it, as well as an error if there were problems
	// marshaling the proto back and forth.
	return safeCopyEntity(e), nil
}

func (m *Manager) updateEntityMeta(e *pb.Entity, newMeta *pb.EntityMeta) error {
	// get the existing metadata
	meta := e.GetMeta()

	// some fields must not be merged in, so we make sure that
	// they're nulled out here
	newMeta.Capabilities = nil
	newMeta.Groups = nil
	newMeta.Keys = nil

	// Merge all changes, and then overwrite the original keys
	// with the ones from newMeta since that is a rewrite style
	// update.
	proto.Merge(meta, newMeta)

	// Save changes
	return m.db.SaveEntity(e)
}

// UpdateEntityMeta drives the internal version by obtaining the
// entity from the database based on the ID.
func (m *Manager) UpdateEntityMeta(entityID string, newMeta *pb.EntityMeta) error {
	e, err := m.db.LoadEntity(entityID)
	if err != nil {
		return err
	}

	return m.updateEntityMeta(e, newMeta)
}

// updateEntityKeys performs an update on keys to allow the client to
// be simpler, and to account for proto.Merge() merging list contents
// rather than overwriting.
func (m *Manager) updateEntityKeys(e *pb.Entity, mode, keyType, key string) ([]string, error) {
	// Normalize the type and the mode
	mode = strings.ToUpper(mode)
	keyType = strings.ToUpper(keyType)

	switch mode {
	case "LIST":
		return e.GetMeta().GetKeys(), nil
	case "ADD":
		e.Meta.Keys = patchStringSlice(e.Meta.Keys, fmt.Sprintf("%s:%s", keyType, key), true, true)
	case "DEL":
		e.Meta.Keys = patchStringSlice(e.Meta.Keys, key, false, false)
	}

	// Save changes
	if err := m.db.SaveEntity(e); err != nil {
		return nil, err
	}
	return nil, nil
}

// UpdateEntityKeys is the exported version of updateEntityKeys
func (m *Manager) UpdateEntityKeys(entityID, mode, keytype, key string) ([]string, error) {
	e, err := m.db.LoadEntity(entityID)
	if err != nil {
		return nil, err
	}

	return m.updateEntityKeys(e, mode, keytype, key)
}

// ManageUntypedEntityMeta handles the things that may be annotated
// onto an entity.  These annotations should be used sparingly as they
// incur a non-trivial lookup cost on the server.
func (m *Manager) ManageUntypedEntityMeta(entityID, mode, key, value string) ([]string, error) {
	// Load Entity
	e, err := m.db.LoadEntity(entityID)
	if err != nil {
		return nil, err
	}

	// Patch the KV slice
	tmp := patchKeyValueSlice(e.GetMeta().UntypedMeta, mode, key, value)

	// If this was a read, bail out now with whatever was read
	if strings.ToUpper(mode) == "READ" {
		return tmp, nil
	}

	// Save changes
	e.Meta.UntypedMeta = tmp
	if err := m.db.SaveEntity(e); err != nil {
		return nil, err
	}
	return nil, nil
}

// setEntityLockState is a convenience function to handle the setting
// of an entity's lock state.
func (m *Manager) setEntityLockState(entityID string, locked bool) error {
	// Load Entity
	e, err := m.db.LoadEntity(entityID)
	if err != nil {
		return err
	}

	// update state
	e.Meta.Locked = &locked

	// Save changes
	if err := m.db.SaveEntity(e); err != nil {
		return err
	}
	return nil
}

// LockEntity allows external callers to lock entities directly.
// Internal users can just set the value directly.
func (m *Manager) LockEntity(entityID string) error {
	return m.setEntityLockState(entityID, true)
}

// UnlockEntity allows external callers to lock entities directly.
// Internal users can just set the value directly.
func (m *Manager) UnlockEntity(entityID string) error {
	return m.setEntityLockState(entityID, false)
}

// safeCopyEntity makes a copy of the entity provided but removes
// fields that are related to security.  This permits the entity that
// is returned to be handed off outside the server.
func safeCopyEntity(e *pb.Entity) *pb.Entity {
	dup := &pb.Entity{}
	proto.Merge(dup, e)

	// Fields for security are nulled out before returning.
	dup.Secret = proto.String("<REDACTED>")

	return dup
}
