package hooks

import (
	"strings"

	"github.com/NetAuth/NetAuth/internal/tree/util"

	pb "github.com/NetAuth/Protocol"
)

type AddEntityUM struct{}

func (*AddEntityUM) Name() string  { return "add-untyped-metadata" }
func (*AddEntityUM) Priority() int { return 50 }
func (*AddEntityUM) Run(e, de *pb.Entity) error {
	for _, m := range de.Meta.UntypedMeta {
		key, value := splitKeyValue(m)
		e.Meta.UntypedMeta = util.PatchKeyValueSlice(e.Meta.UntypedMeta, "UPSERT", key, value)
	}
	return nil
}

type DelFuzzyEntityUM struct{}

func (*DelFuzzyEntityUM) Name() string  { return "del-untyped-metadata-fuzzy" }
func (*DelFuzzyEntityUM) Priority() int { return 50 }
func (*DelFuzzyEntityUM) Run(e, de *pb.Entity) error {
	for _, m := range de.Meta.UntypedMeta {
		key, value := splitKeyValue(m)
		e.Meta.UntypedMeta = util.PatchKeyValueSlice(e.Meta.UntypedMeta, "CLEARFUZZY", key, value)
	}
	return nil
}

type DelExactEntityUM struct{}

func (*DelExactEntityUM) Name() string  { return "del-untyped-metadata-strict" }
func (*DelExactEntityUM) Priority() int { return 50 }
func (*DelExactEntityUM) Run(e, de *pb.Entity) error {
	for _, m := range de.Meta.UntypedMeta {
		key, value := splitKeyValue(m)
		e.Meta.UntypedMeta = util.PatchKeyValueSlice(e.Meta.UntypedMeta, "CLEAREXACT", key, value)
	}
	return nil
}

func splitKeyValue(s string) (string, string) {
	parts := strings.SplitN(s, ":", 2)
	return parts[0], parts[1]
}
