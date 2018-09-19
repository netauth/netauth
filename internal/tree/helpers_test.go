package tree

import (
	"sort"
	"testing"
)

func TestPatchStringSlice(t *testing.T) {
	cases := []struct {
		in         []string
		patch      string
		insert     bool
		matchExact bool
		want       []string
	}{
		{[]string{"foo", "bar", "baz"}, "qux", true, true, []string{"foo", "bar", "baz", "qux"}},
		{[]string{"foo", "bar", "baz", "qux"}, "qux", false, true, []string{"foo", "bar", "baz"}},
		{[]string{"foo", "bar", "baz"}, "qux", false, true, []string{"foo", "bar", "baz"}},
		{[]string{"foo", "bar", "baz"}, "ba", false, false, []string{"foo"}},

		{[]string{"foo", "bar", "baz"}, "foo", true, true, []string{"foo", "bar", "baz"}},
	}

	for i, c := range cases {
		got := patchStringSlice(c.in, c.patch, c.insert, c.matchExact)
		if !slicesAreEqual(got, c.want) {
			t.Errorf("%d: Got %v; Want %v", i, got, c.want)
		}
	}
}

func TestStringMatcher(t *testing.T) {
	cases := []struct {
		str       string
		substr      string
		matchExact bool
		want       bool
	}{
		{"foo", "foo", false, true},
		{"foo", "foo", true, true},
		{"foosball", "foo", false, true},
		{"foosball", "foo", true, false},

		{"foo", "bar", false, false},
		{"foo", "bar", true, false},
	}

	for i, c := range cases {
		if got := stringMatcher(c.str, c.substr, c.matchExact); got != c.want {
			t.Errorf("%d: Got %v; Want %v (%v != %v)", i, got, c.want, c.str, c.substr)
		}
	}
}

func slicesAreEqual(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}

	sort.Slice(left, func(i, j int) bool {
		return left[i] < left[j]
	})

	sort.Slice(right, func(i, j int) bool {
		return right[i] < right[j]
	})

	for i, v := range left {
		if v != right[i] {
			return false
		}
	}
	return true
}

func TestSlicesAreEqual(t *testing.T) {
	cases := []struct {
		left  []string
		right []string
		want  bool
	}{
		{[]string{"foo", "bar", "baz"}, []string{"foo", "bar"}, false},
		{[]string{"foo", "bar", "baz"}, []string{"foo", "bar", "baz"}, true},
		{[]string{"foo", "bar", "baz"}, []string{"foo", "baz", "bar"}, true},
	}

	for i, c := range cases {
		if got := slicesAreEqual(c.left, c.right); got != c.want {
			t.Errorf("%d: Got %v; Want %v", i, got, c.want)
		}
	}
}
