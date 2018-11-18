package hooks

import (
	"strings"
)

func splitKeyValue(s string) (string, string) {
	parts := strings.SplitN(s, ":", 2)
	return parts[0], parts[1]
}
