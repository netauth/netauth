package all

import (
	// The blank import here permits the init() within the bcrypt module
	// to register its implementation to the crypto plugin system.
	_ "github.com/NetAuth/NetAuth/internal/crypto/bcrypt"
)
