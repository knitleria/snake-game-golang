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
