// +build nocrypto

package all

import (
	// The blank import here permits the init() within the
	// nocrypto module to register its implementation to the
	// crypto plugin system.
	_ "github.com/netauth/netauth/internal/crypto/nocrypto"
)
