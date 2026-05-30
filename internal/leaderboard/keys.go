package leaderboard

import (
	"fmt"
	"time"
)

const scoreSortBase = 999999999

func playerPK(playerID string) string {
	return "PLAYER#" + playerID
}

func modeSK(mode string) string {
	return "MODE#" + mode
}

func modeGSIKey(mode string) string {
	return "MODE#" + mode
}

func leaderboardSortKey(score int, updatedAt time.Time, playerID string) string {
	invertedScore := scoreSortBase - score
	return fmt.Sprintf(
		"SCORE#%09d#TIME#%s#PLAYER#%s",
		invertedScore,
		updatedAt.UTC().Format(time.RFC3339Nano),
		playerID,
	)
}
