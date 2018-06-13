package ctl

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/google/subcommands"
)

type ModifyKeysCmd struct {
	ID      string
	keyType string
	key     string
	mode    string
}

func (*ModifyKeysCmd) Name() string     { return "modify-keys" }
func (*ModifyKeysCmd) Synopsis() string { return "Modify stored keys on an entity" }
func (*ModifyKeysCmd) Usage() string {
	return `modify-keys --ID <ID> --mode <ADD|LIST|DEL> --type <type> --key <key>

Modify the stored keys for an entity.  Key type must be specified as a
string of well known type.  This will be an uppercase version of the
key type like 'SSH' or 'GPG'.
`
}

func (p *ModifyKeysCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.ID, "ID", entity, "Entity to act on")
	f.StringVar(&p.keyType, "type", "SSH", "Type of the key")
	f.StringVar(&p.key, "key", "", "Key contents")
	f.StringVar(&p.mode, "mode", "LIST", "Action to perform on keys")
}

func (p *ModifyKeysCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Grab a client
	c, err := getClient()
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	// Get the authorization token
	t, err := c.GetToken(entity, secret)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	keys, err := c.ModifyEntityKeys(t, p.ID, p.mode, p.keyType, p.key)
	if err != nil {
		return subcommands.ExitFailure
	}

	for _, k := range keys {
		parts := strings.Split(k, ":")
		fmt.Printf("Type: %s; Key: %s\n", parts[0], strings.Join(parts[1:], " "))
	}

	return subcommands.ExitSuccess
}
