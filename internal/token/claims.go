package token

// Claims is a type that contains the claims that all tokens shall
// have.  Implementations may embed additional messages, but these
// cliams must exist here.
type Claims struct {
	EntityID     string
	Capabilities []string
	RenewalsLeft int
}

// HasCapability is a convenience function to determine if the
// provided token contains the requested capability.  The capability
// GLOBAL_ROOT will cause the function to return true immediately as
// GLOBAL_ROOT counts for all capabilities.
func (c *Claims) HasCapability(cap string) bool {
	for _, tc := range c.Capabilities {
		if tc == cap || tc == "GLOBAL_ROOT" {
			return true
		}
	}
	return false
}
