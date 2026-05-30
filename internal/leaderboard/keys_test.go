package leaderboard

import (
	"sort"
	"testing"
	"time"
)

func TestLeaderboardSortKey(t *testing.T) {
	now := time.Date(2026, 5, 16, 10, 0, 0, 0, time.UTC)
	keys := []string{
		leaderboardSortKey(10, now, "b"),
		leaderboardSortKey(20, now, "a"),
		leaderboardSortKey(5, now, "c"),
	}
	sort.Strings(keys)

	wantFirst := leaderboardSortKey(20, now, "a")
	if keys[0] != wantFirst {
		t.Fatalf("highest score should sort first: got %s want %s", keys[0], wantFirst)
	}
}
