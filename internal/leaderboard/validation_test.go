package leaderboard

import (
	"errors"
	"testing"
)

func TestValidateSubmitRequest(t *testing.T) {
	valid := SubmitScoreRequest{
		PlayerID:      "11111111111111111111111111111111",
		PlayerName:    "Valeria",
		Score:         10,
		Mode:          "normal",
		Skin:          "normal",
		ClientVersion: "dev",
		DurationMS:    1000,
	}

	tests := []struct {
		name string
		req  SubmitScoreRequest
		err  error
	}{
		{name: "valid", req: valid},
		{name: "bad player id", req: with(valid, func(r *SubmitScoreRequest) { r.PlayerID = "nope" }), err: ErrBadPlayerID},
		{name: "empty name", req: with(valid, func(r *SubmitScoreRequest) { r.PlayerName = " " }), err: ErrEmptyName},
		{name: "long name", req: with(valid, func(r *SubmitScoreRequest) { r.PlayerName = "123456789012345678901" }), err: ErrNameTooLong},
		{name: "bad score low", req: with(valid, func(r *SubmitScoreRequest) { r.Score = 0 }), err: ErrBadScore},
		{name: "bad score high", req: with(valid, func(r *SubmitScoreRequest) { r.Score = 1200 }), err: ErrBadScore},
		{name: "bad mode", req: with(valid, func(r *SubmitScoreRequest) { r.Mode = "hard" }), err: ErrBadMode},
		{name: "bad skin", req: with(valid, func(r *SubmitScoreRequest) { r.Skin = "unknown" }), err: ErrBadSkin},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateSubmitRequest(tt.req, DefaultMaxScore)
			if tt.err == nil {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got.PlayerName == "" {
					t.Fatal("expected normalized player name")
				}
				return
			}
			if !errors.Is(err, tt.err) {
				t.Fatalf("expected %v, got %v", tt.err, err)
			}
		})
	}
}

func with(req SubmitScoreRequest, f func(*SubmitScoreRequest)) SubmitScoreRequest {
	f(&req)
	return req
}
