package mugen_test

import (
	"context"
	"testing"

	"github.com/Japan7/karaberus/server/clients/mugen"
	"github.com/google/uuid"
)

func TestGetKara(t *testing.T) {
	client := mugen.GetClient()
	kid, err := uuid.Parse("a0ac08a9-46b6-4219-8017-a0a9581ef914")
	if err != nil {
		t.Fatal("failed to parse kid")
	}

	resp, err := client.GetKara(context.TODO(), kid)
	if err != nil {
		t.Fatalf("request failed: %s", err)
	}

	if resp.KID != kid {
		t.Fatal("response kid is different from requested kid")
	}
}
