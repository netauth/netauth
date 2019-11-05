// +build MemDB

package all

import (
	// Register the database from init()
	_ "github.com/netauth/netauth/internal/db/memdb"
)
