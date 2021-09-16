package tree

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/netauth/netauth/internal/tree/util"

	"github.com/netauth/netauth/internal/db"

	pb "github.com/netauth/protocol"
)

// CreateEntity creates a new entity given an ID, number, and secret.
// Its not necessary to set the secret upon creation and it can be set
// later.  If not set on creation then the entity will not be usable.
// number must be a unique positive integer.  Because these are
// generally allocated in sequence the special value '-1' may be
// specified which will select the next available number.
func (m *Manager) CreateEntity(ID string, number int32, secret string) error {
	de := &pb.Entity{
		ID:     &ID,
		Number: &number,
		Secret: &secret,
	}

	_, err := m.RunEntityChain("CREATE", de)
	return err
}

// DestroyEntity deletes the named entity.  This function will
// delete the entity in a non-atomic way, but will ensure that the
// entity cannot be authenticated with before returning.  If the named
// ID does not exist the function will return tree.E_NO_ENTITY, in
// all other cases nil is returned.
func (m *Manager) DestroyEntity(ID string) error {
	de := &pb.Entity{
		ID: &ID,
	}

	_, err := m.RunEntityChain("DESTROY", de)
	return err
}

// SetEntityCapability adds a capability to an entry directly.
func (m *Manager) SetEntityCapability(ID string, c string) error {
	capIndex, ok := pb.Capability_value[c]
	if !ok {
		return ErrUnknownCapability
	}
	cap := pb.Capability(capIndex)
	return m.SetEntityCapability2(ID, &cap)
}

// SetEntityCapability2 adds a capability to an entity directly, and
// does so with a strongly typed capability pointer.
func (m *Manager) SetEntityCapability2(ID string, c *pb.Capability) error {
	if c == nil {
		return ErrUnknownCapability
	}

	de := &pb.Entity{
		ID: &ID,
		Meta: &pb.EntityMeta{
			Capabilities: []pb.Capability{*c},
		},
	}

	_, err := m.RunEntityChain("SET-CAPABILITY", de)
	return err
}

// DropEntityCapability adds a capability to an entry directly.
func (m *Manager) DropEntityCapability(ID string, c string) error {
	capIndex, ok := pb.Capability_value[c]
	if !ok {
		return ErrUnknownCapability
	}
	cap := pb.Capability(capIndex)
	return m.DropEntityCapability2(ID, &cap)
}

// DropEntityCapability2 adds a capability to an entity directly, and
// does so with a strongly typed capability pointer.
func (m *Manager) DropEntityCapability2(ID string, c *pb.Capability) error {
	if c == nil {
		return ErrUnknownCapability
	}

	de := &pb.Entity{
		ID: &ID,
		Meta: &pb.EntityMeta{
			Capabilities: []pb.Capability{*c},
		},
	}

	_, err := m.RunEntityChain("DROP-CAPABILITY", de)
	return err
}

// SetSecret sets the secret on a given entity using the
// crypto interface.
func (m *Manager) SetSecret(ID string, secret string) error {
	de := &pb.Entity{
		ID:     &ID,
		Secret: &secret,
	}

	_, err := m.RunEntityChain("SET-SECRET", de)
	return err
}

// ValidateSecret validates the identity of an entity by
// validating the authenticating entity with the secret.
func (m *Manager) ValidateSecret(ID string, secret string) error {
	de := &pb.Entity{
		ID:     &ID,
		Secret: &secret,
	}

	_, err := m.RunEntityChain("VALIDATE-IDENTITY", de)
	return err
}

// FetchEntity returns an entity to the caller after first making a
// safe copy of it to remove secure fields.
func (m *Manager) FetchEntity(ID string) (*pb.Entity, error) {
	de := &pb.Entity{
		ID: &ID,
	}

	e, err := m.RunEntityChain("FETCH", de)
	if err != nil {
		return nil, err
	}

	// The safeCopyEntity will return the entity without secrets
	// in it, as well as an error if there were problems
	// marshaling the proto back and forth.
	return safeCopyEntity(e), nil
}

// UpdateEntityMeta drives the internal version by obtaining the
// entity from the database based on the ID.
func (m *Manager) UpdateEntityMeta(ID string, newMeta *pb.EntityMeta) error {
	de := &pb.Entity{
		ID:   &ID,
		Meta: newMeta,
	}

	_, err := m.RunEntityChain("MERGE-METADATA", de)
	return err
}

// UpdateEntityKeys manages entity public keys.  Additional setup
// occurs to select the correct processing chain based on what action
// was requested.
func (m *Manager) UpdateEntityKeys(ID, mode, keytype, key string) ([]string, error) {
	mode = strings.ToUpper(mode)
	keytype = strings.ToUpper(keytype)

	// Configure request data.
	de := &pb.Entity{
		ID: &ID,
		Meta: &pb.EntityMeta{
			Keys: []string{fmt.Sprintf("%s:%s", keytype, key)},
		},
	}

	// Select chain based on mode, or coerce to 'LIST'
	chain := "FETCH"
	switch mode {
	case "ADD":
		chain = "ADD-KEY"
	case "DROP":
		fallthrough
	case "DEL":
		chain = "DEL-KEY"
	default:
		mode = "LIST"
	}

	// Execute the transaction.
	e, err := m.RunEntityChain(chain, de)
	if err != nil {
		return nil, err
	}

	// If this was just a read request, return the data.
	if mode == "LIST" {
		return e.GetMeta().GetKeys(), nil
	}
	return nil, nil
}

// ManageUntypedEntityMeta handles the things that may be annotated
// onto an entity.  These annotations should be used sparingly as they
// incur a non-trivial lookup cost on the server.
func (m *Manager) ManageUntypedEntityMeta(ID, mode, key, value string) ([]string, error) {
	de := &pb.Entity{
		ID: &ID,
		Meta: &pb.EntityMeta{
			UntypedMeta: []string{fmt.Sprintf("%s:%s", key, value)},
		},
	}

	// Mode switch and select appropriate processor chain.
	chain := "FETCH"
	switch strings.ToUpper(mode) {
	case "UPSERT":
		chain = "UEM-UPSERT"
	case "CLEARFUZZY":
		chain = "UEM-CLEARFUZZY"
	case "CLEAREXACT":
		chain = "UEM-CLEAREXACT"
	default:
		mode = "READ"
	}

	// Process transaction
	e, err := m.RunEntityChain(chain, de)
	if err != nil {
		return nil, err
	}

	// If this was a read, bail out now with whatever was read
	if strings.ToUpper(mode) == "READ" {
		return util.PatchKeyValueSlice(e.GetMeta().GetUntypedMeta(), "READ", key, ""), nil
	}
	return nil, nil
}

// EntityKVAdd handles adding a new key to the KV store for an entity
// identified by ID.  The key must not previously exist.
func (m *Manager) EntityKVAdd(ID string, d []*pb.KVData) error {
	de := &pb.Entity{
		ID: &ID,
		Meta: &pb.EntityMeta{
			KV: d,
		},
	}

	_, err := m.RunEntityChain("KV-ADD", de)
	return err
}

// EntityKVDel handles removing an existing key from the entity
// identified by ID.  An attempt to remove a key that does not exist
// will return an error.
func (m *Manager) EntityKVDel(ID string, d []*pb.KVData) error {
	de := &pb.Entity{
		ID: &ID,
		Meta: &pb.EntityMeta{
			KV: d,
		},
	}

	_, err := m.RunEntityChain("KV-DEL", de)
	return err
}

// EntityKVReplace handles replacing an existing key on the entity
// identified by ID.  An attempt to replace a key that does not exist
// will return an error.
func (m *Manager) EntityKVReplace(ID string, d []*pb.KVData) error {
	de := &pb.Entity{
		ID: &ID,
		Meta: &pb.EntityMeta{
			KV: d,
		},
	}

	_, err := m.RunEntityChain("KV-REPLACE", de)
	return err
}

// EntityKVGet returns a selected key or keys to the caller.
func (m *Manager) EntityKVGet(ID string, keys []*pb.KVData) ([]*pb.KVData, error) {
	e, err := m.FetchEntity(ID)
	if err != nil {
		return nil, err
	}

	if len(keys) == 1 && keys[0].GetKey() == "*" {
		// In the special case of a single star as the key
		// return the entire keyspace
		return e.GetMeta().GetKV(), nil
	}

	out := []*pb.KVData{}
	for _, haystack := range e.GetMeta().GetKV() {
		for _, needle := range keys {
			if haystack.GetKey() != needle.GetKey() {
				continue
			}
			out = append(out, haystack)
		}
	}
	if len(out) == 0 {
		return nil, ErrNoSuchKey
	}
	return out, nil
}

// LockEntity allows external callers to lock entities directly.
// Internal users can just set the value directly.
func (m *Manager) LockEntity(ID string) error {
	de := &pb.Entity{
		ID: &ID,
	}

	_, err := m.RunEntityChain("LOCK", de)
	return err
}

// UnlockEntity allows external callers to lock entities directly.
// Internal users can just set the value directly.
func (m *Manager) UnlockEntity(ID string) error {
	de := &pb.Entity{
		ID: &ID,
	}

	_, err := m.RunEntityChain("UNLOCK", de)
	return err
}

func (m *Manager) entityResolverCallback(e db.Event) {
	switch e.Type {
	case db.EventEntityCreate:
		fallthrough
	case db.EventEntityUpdate:
		ent, err := m.db.LoadEntity(e.PK)
		if err != nil {
			m.log.Warn("Unchecked load error in entityResolverCallback", "error", err)
			return
		}
		m.resolver.SyncDirectGroups(ent.GetID(), ent.GetMeta().GetGroups())
	case db.EventEntityDestroy:
		m.resolver.RemoveEntity(e.PK)
	default:
		return
	}
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
