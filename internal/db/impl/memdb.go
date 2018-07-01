// +build MemDB

package impl

import (
	// Register the database from init()
	_ "github.com/NetAuth/NetAuth/internal/db/impl/memdb"
)
