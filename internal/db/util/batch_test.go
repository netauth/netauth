package util

import (
	"testing"
)

func TestLoadEntityBatch(t *testing.T) {
	s, err := LoadEntityBatch([]string{"entity1", "entity2"}, okEntityLoader)
	if err != nil {
		t.Fatal(err)
	}
	if len(s) != 2 {
		t.Error("Batch load returned the wrong number of results")
	}
}

func TestLoadEntityBatchBadLoader(t *testing.T) {
	s, err := LoadEntityBatch([]string{"entity1", "entity2"}, evilEntityLoader)
	if s != nil || err == nil {
		t.Error("Batch operation returned results with a non-nil error")
	}
}

func TestLoadGroupBatch(t *testing.T) {
	s, err := LoadGroupBatch([]string{"group1", "group2"}, okGroupLoader)
	if err != nil {
		t.Fatal(err)
	}
	if len(s) != 2 {
		t.Error("Batch load returned the wrong number of results")
	}
}

func TestLoadGroupBatchBadLoader(t *testing.T) {
	s, err := LoadGroupBatch([]string{"group1", "group2"}, evilGroupLoader)
	if s != nil || err == nil {
		t.Error("Batch operation returned results with a non-nil error")
	}
}
