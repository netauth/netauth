package tree

import (
	"fmt"
	"regexp"
	"sort"
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

// patchKeyValueSlice patches slices that use key/value pairs.  Its
// designed with more advanced functionality around exact key
// matching, fuzzy and exact clearing, and OpenLDAP-style Z-Ordering.
func patchKeyValueSlice(slice []string, mode, key, value string) []string {
	mode = strings.ToUpper(mode)

	switch mode {
	case "UPSERT":
		var newSlice []string
		inserted := false
		for _, s := range slice {
			parts := strings.Split(s, ":")
			if parts[0] == key {
				newSlice = append(newSlice, fmt.Sprintf("%s:%s", key, value))
				inserted = true
			} else {
				newSlice = append(newSlice, s)
			}
		}
		if !inserted {
			newSlice = append(newSlice, fmt.Sprintf("%s:%s", key, value))
		}

		newSlice = dedupStringSlice(newSlice)
		sort.Strings(newSlice)
		return newSlice
	case "CLEARFUZZY":
		var newSlice []string
		// Iterate over the keys, performing matching after
		// discarding the pattern {\d+}$ to permit OpenLDAP
		// style Z-Ordering of values.
		re := regexp.MustCompile("{\\d+}$")
		strippedK := re.ReplaceAllString(key, "")
		for _, kv := range slice {
			parts := strings.Split(kv, ":")
			if re.ReplaceAllString(parts[0], "") != strippedK {
				newSlice = append(newSlice, kv)
			}
		}

		newSlice = dedupStringSlice(newSlice)
		sort.Strings(newSlice)
		return newSlice
	case "CLEAREXACT":
		var newSlice []string

		for _, kv := range slice {
			parts := strings.Split(kv, ":")
			if parts[0] != key {
				newSlice = append(newSlice, kv)
			}
		}

		newSlice = dedupStringSlice(newSlice)
		sort.Strings(newSlice)
		return newSlice
	case "READ":
		// Special key '*' returns all values.
		out := []string{}
		if key == "*" {
			out = slice
			sort.Strings(out)
			return out
		}

		// Iterate over the keys, performing matching after
		// discarding the pattern {\d+}$ to permit OpenLDAP
		// style Z-Ordering of values.
		re := regexp.MustCompile("{\\d+}$")
		strippedK := re.ReplaceAllString(key, "")
		for _, kv := range slice {
			parts := strings.Split(kv, ":")
			if re.ReplaceAllString(parts[0], "") == strippedK {
				out = append(out, kv)
			}
		}

		// Sort to ensure any Z-Ordering that may be present
		sort.Strings(out)

		return out
	}

	// We shouldn't be here, but just in case return the original
	// slice unmodified for safety.
	return slice
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
