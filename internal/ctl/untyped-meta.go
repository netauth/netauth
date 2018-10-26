package ctl

import (
	"context"
	"flag"
	"fmt"
	"sort"

	"github.com/google/subcommands"
)

// UntypedMetaCmd requests the server to modify the untyped metadata
// on groups and entities.
type UntypedMetaCmd struct {
	key        string
	value      string
	entityID   string
	groupName  string
	read       bool
	clearexact bool
	clearfuzzy bool
	upsert     bool
}

// Name of this cmdlet is 'untyped-meta'
func (*UntypedMetaCmd) Name() string { return "untyped-meta" }

// Synopsis returns short-form usage information.
func (*UntypedMetaCmd) Synopsis() string { return "Manage untyped meta on entities and groups" }

// Usage returns long-form usage information.
func (*UntypedMetaCmd) Usage() string {
	return `untyped-meta --<entity|group> --<read|clearfuzzy|clearexact|upsert> --key <key> --value [value]

Manage untyped metadata associated with groups and entities.
Clearfuzzy strips ordering values from keys before clearing, whereas
clearexact removes a key with an exact match.  Upsert either inserts a
new value, or updates an existing one, and tries to do the right thing
either way.  Key and Value are always strings, key may not contain the
reserved rune ':'.
`
}

// SetFlags sets the cmdlet specific flags.
func (p *UntypedMetaCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.entityID, "entity", "", "ID for the entity to modify")
	f.StringVar(&p.groupName, "group", "", "Name for the group to modify")
	f.StringVar(&p.key, "key", "", "Key to act on")
	f.StringVar(&p.value, "value", "", "Value to act with")
	f.BoolVar(&p.read, "read", false, "Read the value specified by 'key' including ordered results")
	f.BoolVar(&p.clearexact, "clear-exact", false, "Clear keys with exact matching")
	f.BoolVar(&p.clearfuzzy, "clear-fuzzy", false, "Clear keys with partial matching")
	f.BoolVar(&p.upsert, "upsert", false, "Insert/Update a value for a particular key")
}

// Execute runs the cmdlet.
func (p *UntypedMetaCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Only act on one type at a time
	if (p.entityID != "" && p.groupName != "") || (p.entityID == "" && p.groupName == "") {
		fmt.Println("Exactly one of --entity or --group must be specified")
		return subcommands.ExitFailure
	}

	// Only accept one mode operation
	modeOpts := []bool{p.read, p.clearexact, p.clearfuzzy, p.upsert}
	tAny := false
	tTwo := false
	for m := range modeOpts {
		tTwo = (tAny && modeOpts[m])
		tAny = (tAny || modeOpts[m])
	}
	if tTwo {
		fmt.Println("Exactly one of --read, --upsert, --clear-fuzzy, or --clear-exact must be specified.")
		return subcommands.ExitFailure
	}

	// Set the mode
	mode := ""
	if p.read {
		mode = "READ"
	} else if p.clearexact {
		mode = "CLEAREXACT"
	} else if p.clearfuzzy {
		mode = "CLEARFUZZY"
	} else if p.upsert {
		mode = "UPSERT"
	}

	// Grab a client
	c, err := getClient()
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	// Get the authorization token if needed
	t := ""
	if mode != "READ" {
		t, err = getToken(c, getEntity())
		if err != nil {
			fmt.Println(err)
			return subcommands.ExitFailure
		}
	}
	// Query the server
	result := make(map[string]string)

	if p.entityID != "" {
		result, err = c.ModifyUntypedEntityMeta(t, p.entityID, mode, p.key, p.value)
	} else if p.groupName != "" {
		result, err = c.ModifyUntypedGroupMeta(t, p.groupName, mode, p.key, p.value)
	}

	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	out := []string{}
	for k, v := range result {
		out = append(out, fmt.Sprintf("%s: %s", k, v))
	}
	sort.Strings(out)

	for _, l := range out {
		fmt.Println(l)
	}

	return subcommands.ExitSuccess
}
