package ctl

import (
	"context"
	"flag"
	"fmt"

	pb "github.com/NetAuth/Protocol"

	"github.com/google/subcommands"
)

// LockEntityCmd requests the server to lock or unlock an entity.
type LockEntityCmd struct {
	entityID string
	lock     bool
	unlock   bool
}

// Name of this cmdlet is 'lock-entity'
func (*LockEntityCmd) Name() string { return "lock-entity" }

// Synopsis returns short-form usage information.
func (*LockEntityCmd) Synopsis() string { return "Modify the lock state of an entity" }

// Usage returns long-form usage information.
func (*LockEntityCmd) Usage() string {
	return `lock-entity --<lock|unlock> --entity <entityID>

Lock or unlock an entity.  Each action requires a specific capability.
`
}

// SetFlags sets the cmdlet specific flags.
func (p *LockEntityCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.entityID, "entity", getEntity(), "ID for the entity to modify")
	f.BoolVar(&p.lock, "lock", false, "Lock the named entity")
	f.BoolVar(&p.unlock, "unlock", false, "Unlock the named entity")
}

// Execute runs the cmdlet.
func (p *LockEntityCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
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

	var result *pb.SimpleResult
	if p.lock && !p.unlock {
		result, err = c.LockEntity(t, p.entityID)
	} else if !p.lock && p.unlock {
		result, err = c.UnlockEntity(t, p.entityID)
	} else {
		fmt.Println("Exactly one of '--lock' or '--unlock' must be specified")
		return subcommands.ExitFailure
	}

	// Parse the result
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}
	if result.GetMsg() != "" {
		fmt.Println(result.GetMsg())
	}

	return subcommands.ExitSuccess
}
