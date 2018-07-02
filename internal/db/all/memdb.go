// +build MemDB

package all

import (
	// Register the database from init()
	_ "github.com/NetAuth/NetAuth/internal/db/memdb"
)
