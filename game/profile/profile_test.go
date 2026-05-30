package profile

import (
	"strings"
	"testing"
)

func TestNormalizePlayerName(t *testing.T) {
	t.Parallel()

	name, err := NormalizePlayerName("  Valeria  ")
	if err != nil {
		t.Fatalf("NormalizePlayerName returned error: %v", err)
	}
	if name != "Valeria" {
		t.Fatalf("NormalizePlayerName = %q, want %q", name, "Valeria")
	}
}

func TestNormalizePlayerNameRejectsInvalidValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want error
	}{
		{name: "   ", want: ErrEmptyName},
		{name: strings.Repeat("a", MaxPlayerNameRunes+1), want: ErrNameTooLong},
		{name: "bad\nname", want: ErrBadName},
	}

	for _, tt := range tests {
		_, err := NormalizePlayerName(tt.name)
		if err != tt.want {
			t.Fatalf("NormalizePlayerName(%q) error = %v, want %v", tt.name, err, tt.want)
		}
	}
}

func TestEnsurePlayerID(t *testing.T) {
	t.Parallel()

	config := Config{PlayerName: "Valeria"}
	if err := EnsurePlayerID(&config); err != nil {
		t.Fatalf("EnsurePlayerID: %v", err)
	}
	if len(config.PlayerID) != 32 {
		t.Fatalf("expected 32-char player id, got %q", config.PlayerID)
	}
}

func TestEnsurePlayerIDPreservesExisting(t *testing.T) {
	t.Parallel()

	const existing = "abcdefabcdefabcdefabcdefabcdefab"
	config := Config{PlayerID: existing}
	if err := EnsurePlayerID(&config); err != nil {
		t.Fatalf("EnsurePlayerID: %v", err)
	}
	if config.PlayerID != existing {
		t.Fatalf("EnsurePlayerID overrode existing id: got %q want %q", config.PlayerID, existing)
	}
}

func TestNormalizePlayerID(t *testing.T) {
	t.Parallel()

	id, err := NormalizePlayerID("ABCDEFABCDEFABCDEFABCDEFABCDEFAB")
	if err != nil {
		t.Fatalf("NormalizePlayerID: %v", err)
	}
	if id != "abcdefabcdefabcdefabcdefabcdefab" {
		t.Fatalf("unexpected normalized id: %s", id)
	}
}

func TestNormalizePlayerIDRejectsBad(t *testing.T) {
	t.Parallel()

	if _, err := NormalizePlayerID("not-a-hex-id"); err != ErrBadPlayerID {
		t.Fatalf("expected ErrBadPlayerID, got %v", err)
	}
}
