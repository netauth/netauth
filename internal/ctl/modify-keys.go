package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
)

// ModifyKeysCmd adds, removes, and lists the keys visible on an entity.
type ModifyKeysCmd struct {
	entityID string
	keyType  string
	key      string
	mode     string
}

// Name of this cmdlet is 'modify-keys'
func (*ModifyKeysCmd) Name() string { return "modify-keys" }

// Synopsis returns the short-form usage information.
func (*ModifyKeysCmd) Synopsis() string { return "Modify stored keys on an entity" }

// Usage returns the long-form usage information.
func (*ModifyKeysCmd) Usage() string {
	return `modify-keys --entity <ID> --mode <ADD|LIST|DEL> --type <type> --key <key>

Modify the stored keys for an entity.  Key type must be specified as a
string of well known type.  This will be an uppercase version of the
key type like 'SSH' or 'GPG'.
`
}

// SetFlags sets the cmdlet specific flags.
func (p *ModifyKeysCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.entityID, "entity", getEntity(), "Entity to act on")
	f.StringVar(&p.keyType, "type", "SSH", "Type of the key")
	f.StringVar(&p.key, "key", "", "Key contents")
	f.StringVar(&p.mode, "mode", "LIST", "Action to perform on keys")
}

// Execute runs the cmdlet.
func (p *ModifyKeysCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Grab a client
	c, err := getClient()
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	// Get the authorization token
	t, err := getToken(c, getEntity())
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	keys, err := c.ModifyEntityKeys(t, p.entityID, p.mode, p.keyType, p.key)
	if err != nil {
		return subcommands.ExitFailure
	}

	for _, k := range keys {
		fmt.Printf("Type: %s; Key: %s\n", p.keyType, k)
	}

	return subcommands.ExitSuccess
}
