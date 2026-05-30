package game

import (
	"context"
	"strings"
	"time"

	"snake_golang/assets/mods"
	"snake_golang/assets/skins"
	gamelb "snake_golang/game/leaderboard"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type scoreSubmitState int

const (
	scoreSubmitIdle scoreSubmitState = iota
	scoreSubmitDisabled
	scoreSubmitSubmitting
	scoreSubmitSucceeded
	scoreSubmitFailed
)

type scoreSubmitResult struct {
	response gamelb.SubmitScoreResponse
	err      error
}

type leaderboardFetchResult struct {
	response gamelb.LeaderboardResponse
	err      error
}

func (s *Screen) resetScoreSubmitState() {
	s.scoreSubmitCh = nil
	s.scoreSubmitState = scoreSubmitIdle
	s.scoreSubmitRank = 0
	s.scoreSubmitImproved = false
	s.scoreSubmitError = ""
}

func (s *Screen) pollLeaderboardAsync() {
	if s.scoreSubmitCh != nil {
		select {
		case result := <-s.scoreSubmitCh:
			s.scoreSubmitCh = nil
			if result.err != nil {
				s.scoreSubmitState = scoreSubmitFailed
				s.scoreSubmitError = result.err.Error()
			} else {
				s.scoreSubmitState = scoreSubmitSucceeded
				s.scoreSubmitRank = result.response.Rank
				s.scoreSubmitImproved = result.response.Improved
			}
		default:
		}
	}

	if s.leaderboardCh != nil {
		select {
		case result := <-s.leaderboardCh:
			s.leaderboardCh = nil
			s.leaderboardLoading = false
			if result.err != nil {
				s.leaderboardError = result.err.Error()
				s.leaderboardEntries = nil
			} else {
				s.leaderboardError = ""
				s.leaderboardEntries = result.response.Entries
			}
		default:
		}
	}
}

func (s *Screen) onGameOver() {
	if s.World == nil {
		return
	}
	if s.LeaderboardClient == nil || !s.LeaderboardClient.Enabled() {
		s.scoreSubmitState = scoreSubmitDisabled
		return
	}
	if s.scoreSubmitState == scoreSubmitSubmitting || s.scoreSubmitState == scoreSubmitSucceeded {
		return
	}
	if strings.TrimSpace(s.PlayerName) == "" || strings.TrimSpace(s.PlayerID) == "" {
		s.scoreSubmitState = scoreSubmitDisabled
		return
	}
	if s.World.Score <= 0 {
		s.scoreSubmitState = scoreSubmitDisabled
		return
	}

	durationMS := int64(0)
	if !s.World.StartedAt.IsZero() {
		durationMS = time.Since(s.World.StartedAt).Milliseconds()
	}

	req := gamelb.SubmitScoreRequest{
		PlayerID:      s.PlayerID,
		PlayerName:    s.PlayerName,
		Score:         s.World.Score,
		Mode:          mods.Current().ID(),
		Skin:          skins.Current().ID(),
		ClientVersion: "dev",
		DurationMS:    durationMS,
	}

	ch := make(chan scoreSubmitResult, 1)
	s.scoreSubmitCh = ch
	s.scoreSubmitState = scoreSubmitSubmitting
	s.scoreSubmitError = ""

	client := s.LeaderboardClient
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		response, err := client.SubmitScore(ctx, req)
		ch <- scoreSubmitResult{response: response, err: err}
	}()
}

func (s *Screen) OpenLeaderboard() {
	if s.World == nil {
		return
	}
	s.World.State = StateLeaderboard
	s.fetchLeaderboard()
}

func (s *Screen) UpdateLeaderboard() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		s.World.State = StateMenu
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		s.fetchLeaderboard()
	}
	return nil
}

func (s *Screen) fetchLeaderboard() {
	if s.LeaderboardClient == nil || !s.LeaderboardClient.Enabled() {
		s.leaderboardLoading = false
		s.leaderboardError = "disabled"
		s.leaderboardEntries = nil
		return
	}
	if s.leaderboardLoading {
		return
	}

	ch := make(chan leaderboardFetchResult, 1)
	s.leaderboardCh = ch
	s.leaderboardLoading = true
	s.leaderboardError = ""

	client := s.LeaderboardClient
	mode := mods.Current().ID()
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		response, err := client.FetchLeaderboard(ctx, mode, 20)
		ch <- leaderboardFetchResult{response: response, err: err}
	}()
}
