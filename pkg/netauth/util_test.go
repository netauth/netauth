package netauth

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestAuthorize(t *testing.T) {
	ctx := Authorize(context.Background(), "test-token")

	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		t.Fatal("Bad metadata")
	}

	res := md.Get("authorization")
	if res[0] != "test-token" {
		t.Error("Authorization was not correctly attached")
	}
}

func TestParseKV(t *testing.T) {
	kv1 := []string{
		"key{1}:value1",
		"key{3}:value3",
		"key{2}:value2",
		"k{2}:v2",
		"k{1}:v1",
	}

	res := parseKV(kv1)

	if len(res["key"]) != 3 {
		t.Errorf("key output is the wrong length: %v", res["key"])
	}

	if res["k"][1] != "v2" {
		t.Errorf("k does not contain the correct sorted value!: %v", res["k"])
	}
}
