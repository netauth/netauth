package netauth

import (
	"context"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"google.golang.org/grpc/metadata"
)

var (
	kvIndexRegexp = regexp.MustCompile("\\{(\\d+)\\}")
)

// Authorize attaches a token to a provided context, returning a new
// context that is authorized to make calls with the provided token.
func Authorize(ctx context.Context, token string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", token)
}

// parseKV turns an unsorted list of strings into a map of key to
// sorted values.
func parseKV(in []string) map[string][]string {
	type sortableValue struct {
		index int
		value string
	}

	// Get a map of keys from the input list
	typedKV := make(map[string][]sortableValue)
	for _, p := range in {
		parts := strings.SplitN(p, ":", 2)

		match := kvIndexRegexp.FindStringSubmatch(parts[0])
		idx := -1
		if len(match) > 0 {
			// We can throw the error away here because if
			// there is one, the sentinel value of -1 will
			// simply not be cleared.
			idx, _ = strconv.Atoi(match[1])
		}
		strippedK := kvIndexRegexp.ReplaceAllString(parts[0], "")

		typedKV[strippedK] = append(typedKV[strippedK], sortableValue{idx, parts[1]})
	}

	// Sort the values in each key based on the index value.
	for keyset := range typedKV {
		sort.Slice(typedKV[keyset], func(i, j int) bool {
			return typedKV[keyset][i].index < typedKV[keyset][j].index
		})
	}

	// Convert back down to simple types.
	out := make(map[string][]string)
	for k, v := range typedKV {
		for _, val := range v {
			out[k] = append(out[k], val.value)
		}
	}
	return out
}

func (c *Client) appendMetadata(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx,
		"client-name", c.clientName,
		"service-name", c.serviceName,
	)
}
