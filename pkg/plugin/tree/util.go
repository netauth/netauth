package tree

import (
	"github.com/NetAuth/NetAuth/internal/tree/util"
)

// PatchKeyValueSlice is a helper function that makes working with
// key/value pairs easier.
func PatchKeyValueSlice(slice []string, mode, key, value string) []string {
	return util.PatchKeyValueSlice(slice, mode, key, value)
}
