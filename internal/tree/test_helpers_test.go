package tree

import (
	"sort"
	"testing"

	"github.com/netauth/netauth/internal/crypto/nocrypto"
	"github.com/netauth/netauth/internal/db/memdb"
)

func getNewEntityManager(t *testing.T) *Manager {
	db, err := memdb.New()
	if err != nil {
		t.Fatal(err)
	}

	crypto, err := nocrypto.New()
	if err != nil {
		t.Fatal(err)
	}

	m, err := New(db, crypto)
	if err != nil {
		t.Fatal(err)
	}
	return m
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
