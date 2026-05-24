package main

import (
	"testing"

	"github.com/google/uuid"
)

func TestGenerateUuidAndCopyToClipboard(t *testing.T) {
	t.Run("valid UUID when clipboard disabled", func(t *testing.T) {
		u := generateUuidAndCopyToClipboard(nil, false)
		if u == "" {
			t.Fatal("expected non-empty uuid")
		}
		if _, err := uuid.Parse(u); err != nil {
			t.Fatalf("invalid uuid returned: %v", err)
		}
	})

	t.Run("unique UUIDs across calls", func(t *testing.T) {
		u1 := generateUuidAndCopyToClipboard(nil, false)
		u2 := generateUuidAndCopyToClipboard(nil, false)
		if u1 == u2 {
			t.Fatalf("expected different uuids, got identical: %s", u1)
		}
	})
}
