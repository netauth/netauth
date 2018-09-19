package tree

import (
	"strings"
)

// patchStringSlice patches a string into or out of a slice of other
// strings.  It also ensures that the strings are unique within the
// slice.  When insert is false, the action of the function is to
// remove the provided patch string from the input slice.
func patchStringSlice(in []string, patch string, insert bool, matchExact bool) []string {
	var retSlice []string
	inserted := false
	for _, s := range in {
		matched := stringMatcher(s, patch, matchExact)
		if matched && !insert {
			// Continue without copying, patch out of the
			// list
			continue
		} else if matched && insert {
			// Note that the string was already there in
			// the list
			inserted = true
		}
		retSlice = append(retSlice, s)
	}
	if !inserted && insert {
		retSlice = append(retSlice, patch)
	}

	// We return the dedup'd version rather than the normal as the
	// above process doesn't remove dups that may have gotten into
	// the slice in previous versions of NetAuth.
	return dedupStringSlice(retSlice)
}

// stringMatcher solves the problem introduced above of possibly
// matching with exact string matching or partial string matching.
func stringMatcher(test, match string, matchExact bool) bool {
	if matchExact {
		// Looking for an exact match, case sensitive
		return test == match
	}
	// We can match substrings, so we use
	// strings.Contains()
	return strings.Contains(test, match)
}

// dedupStringSlice converts to a map and then back to a string slice
// to dedup strings by exact matches.
func dedupStringSlice(list []string) []string {
	tmp := make(map[string]int)

	// Into a map
	for _, s := range list {
		tmp[s]++
	}

	// Out of the map
	var out []string
	for k := range tmp {
		out = append(out, k)
	}
	return out
}
