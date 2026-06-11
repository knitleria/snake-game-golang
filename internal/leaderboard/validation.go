package leaderboard

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

const (
	DefaultMaxScore     = 1199
	MaxPlayerNameRunes  = 20
	MaxClientVersionLen = 32
	MaxDurationMS       = int64(24 * 60 * 60 * 1000)
)

var (
	playerIDHexPattern = regexp.MustCompile(`^[a-f0-9]{32}$`)
	ErrBadPlayerID     = errors.New("bad player_id")
	ErrEmptyName       = errors.New("player_name is empty")
	ErrNameTooLong     = errors.New("player_name is too long")
	ErrBadName         = errors.New("player_name contains unsupported characters")
	ErrBadScore        = errors.New("score is out of range")
	ErrBadMode         = errors.New("unknown mode")
	ErrBadSkin         = errors.New("unknown skin")
	ErrBadClientVer    = errors.New("client_version is too long")
	ErrBadDuration     = errors.New("duration_ms is out of range")
	ErrClientTooOld    = errors.New("client version is outdated, please update the game")
)

func ValidateSubmitRequest(req SubmitScoreRequest, maxScore int) (SubmitScoreRequest, error) {
	req.PlayerID = strings.ToLower(strings.TrimSpace(req.PlayerID))
	if !playerIDHexPattern.MatchString(req.PlayerID) {
		return req, ErrBadPlayerID
	}

	name, err := normalizePlayerName(req.PlayerName)
	if err != nil {
		return req, err
	}
	req.PlayerName = name

	if maxScore <= 0 {
		maxScore = DefaultMaxScore
	}
	if req.Score <= 0 || req.Score > maxScore {
		return req, ErrBadScore
	}

	req.Mode = strings.TrimSpace(req.Mode)
	if !validMode(req.Mode) {
		return req, ErrBadMode
	}

	req.Skin = strings.TrimSpace(req.Skin)
	if !validSkin(req.Skin) {
		return req, ErrBadSkin
	}

	req.ClientVersion = strings.TrimSpace(req.ClientVersion)
	if len(req.ClientVersion) > MaxClientVersionLen {
		return req, ErrBadClientVer
	}

	if req.DurationMS < 0 || req.DurationMS > MaxDurationMS {
		return req, ErrBadDuration
	}

	return req, nil
}

func normalizePlayerName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", ErrEmptyName
	}
	runes := []rune(name)
	if len(runes) > MaxPlayerNameRunes {
		return "", ErrNameTooLong
	}
	for _, r := range runes {
		if unicode.IsControl(r) {
			return "", ErrBadName
		}
	}
	return name, nil
}

func validMode(mode string) bool {
	switch mode {
	case "normal", "defaltyk":
		return true
	default:
		return false
	}
}

func validSkin(skin string) bool {

	switch skin {
	case "normal", "halfup", "rantlol":
		return true
	default:
		return false
	}
}
