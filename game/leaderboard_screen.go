package game

import (
	"fmt"
	"image/color"

	"snake_golang/assets"
	"snake_golang/assets/mods"
	"snake_golang/game/i18n"

	"github.com/hajimehoshi/ebiten/v2"
	etxt "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func DrawGameOverSubmitStatus(dst *ebiten.Image, face *etxt.GoTextFaceSource, s *Screen) {
	if face == nil || s == nil {
		return
	}

	text := ""
	clr := color.RGBA{R: 230, G: 235, B: 240, A: 255}

	switch s.scoreSubmitState {
	case scoreSubmitSubmitting:
		text = i18n.T("game.score_submitting")
	case scoreSubmitSucceeded:
		if s.scoreSubmitImproved {
			text = fmt.Sprintf(i18n.T("game.score_submitted"), s.scoreSubmitRank)
		} else {
			text = fmt.Sprintf(i18n.T("game.score_best_unchanged"), s.scoreSubmitRank)
		}
	case scoreSubmitFailed:
		text = i18n.T("game.score_submit_failed")
		clr = color.RGBA{R: 255, G: 150, B: 150, A: 255}
	case scoreSubmitDisabled:
		return
	default:
		return
	}

	size := UiLabelFontSize() * 0.34
	if size < 10 {
		size = 10
	}
	if size > 18 {
		size = 18
	}
	size = FitFontSize(face, text, size, float64(ScreenWidth)*0.9, 8)

	fo := &etxt.GoTextFace{Source: face, Size: size}
	op := &etxt.DrawOptions{}
	op.GeoM.Translate(float64(ScreenWidth)/2, float64(ScreenHeight)*0.58)
	op.PrimaryAlign = etxt.AlignCenter
	op.SecondaryAlign = etxt.AlignCenter
	op.ColorScale.ScaleWithColor(clr)
	etxt.Draw(dst, text, fo, op)
}

func DrawLeaderboardScreen(dst *ebiten.Image, face *etxt.GoTextFaceSource, s *Screen) {
	drawLeaderboardBackground(dst)
	if face == nil || s == nil {
		return
	}

	sw, sh := dst.Bounds().Dx(), dst.Bounds().Dy()
	title := i18n.T("leaderboard.title")
	mode := i18n.T(mods.Current().LabelKey())

	drawLeaderboardText(dst, face, title, float64(sw)/2, float64(sh)*0.14, leaderboardTitleFontSize(sw, sh), color.White, etxt.AlignCenter)
	drawLeaderboardText(dst, face, mode, float64(sw)/2, float64(sh)*0.21, leaderboardHintFontSize(sw, sh), color.RGBA{R: 220, G: 225, B: 230, A: 255}, etxt.AlignCenter)

	topY := int(float64(sh) * 0.29)
	rowH := leaderboardRowHeight(sw, sh)
	left := int(float64(sw) * 0.18)
	right := int(float64(sw) * 0.82)

	headerColor := color.RGBA{R: 205, G: 215, B: 230, A: 255}
	drawLeaderboardText(dst, face, i18n.T("leaderboard.rank_header"), float64(left), float64(topY), leaderboardHeaderFontSize(sw, sh), headerColor, etxt.AlignStart)
	drawLeaderboardText(dst, face, i18n.T("leaderboard.name_header"), float64(left+70), float64(topY), leaderboardHeaderFontSize(sw, sh), headerColor, etxt.AlignStart)
	drawLeaderboardText(dst, face, i18n.T("leaderboard.score_header"), float64(right), float64(topY), leaderboardHeaderFontSize(sw, sh), headerColor, etxt.AlignEnd)

	if s.leaderboardLoading {
		drawLeaderboardText(dst, face, i18n.T("leaderboard.loading"), float64(sw)/2, float64(sh)*0.48, leaderboardHintFontSize(sw, sh), color.White, etxt.AlignCenter)
	} else if s.leaderboardError != "" {
		drawLeaderboardText(dst, face, i18n.T("leaderboard.error"), float64(sw)/2, float64(sh)*0.48, leaderboardHintFontSize(sw, sh), color.RGBA{R: 255, G: 150, B: 150, A: 255}, etxt.AlignCenter)
	} else if len(s.leaderboardEntries) == 0 {
		drawLeaderboardText(dst, face, i18n.T("leaderboard.empty"), float64(sw)/2, float64(sh)*0.48, leaderboardHintFontSize(sw, sh), color.White, etxt.AlignCenter)
	} else {
		y := topY + rowH
		for _, entry := range s.leaderboardEntries {
			if y > sh-int(float64(sh)*0.14) {
				break
			}
			rowColor := color.RGBA{R: 245, G: 247, B: 250, A: 255}
			drawLeaderboardText(dst, face, fmt.Sprintf("#%d", entry.Rank), float64(left), float64(y), leaderboardRowFontSize(sw, sh), rowColor, etxt.AlignStart)
			drawLeaderboardText(dst, face, entry.PlayerName, float64(left+70), float64(y), leaderboardRowFontSize(sw, sh), rowColor, etxt.AlignStart)
			drawLeaderboardText(dst, face, fmt.Sprintf("%d", entry.Score), float64(right), float64(y), leaderboardRowFontSize(sw, sh), rowColor, etxt.AlignEnd)
			y += rowH
		}
	}

	drawLeaderboardText(dst, face, i18n.T("leaderboard.refresh_hint"), float64(sw)/2, float64(sh)*0.92, leaderboardHintFontSize(sw, sh), color.RGBA{R: 220, G: 225, B: 230, A: 230}, etxt.AlignCenter)
}

func drawLeaderboardBackground(dst *ebiten.Image) {
	active := assets.Current()
	if active == nil {
		return
	}
	dw, dh := dst.Bounds().Dx(), dst.Bounds().Dy()
	bg := active.MenuBackground
	if bg == nil {
		bg = active.Background
	}
	if bg != nil {
		op := &ebiten.DrawImageOptions{}
		bw := float64(bg.Bounds().Dx())
		bh := float64(bg.Bounds().Dy())
		if bw > 0 && bh > 0 {
			op.GeoM.Scale(float64(dw)/bw, float64(dh)/bh)
		}
		dst.DrawImage(bg, op)
	}
	vector.FillRect(dst, 0, 0, float32(dw), float32(dh), color.RGBA{R: 18, G: 20, B: 28, A: 170}, true)
}

func drawLeaderboardText(dst *ebiten.Image, face *etxt.GoTextFaceSource, text string, x, y, size float64, clr color.Color, align etxt.Align) {
	if text == "" {
		return
	}
	fo := &etxt.GoTextFace{Source: face, Size: size}
	op := &etxt.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.PrimaryAlign = align
	op.SecondaryAlign = etxt.AlignCenter
	op.ColorScale.ScaleWithColor(clr)
	etxt.Draw(dst, text, fo, op)
}

func leaderboardTitleFontSize(sw, sh int) float64 {
	size := float64(min(sw, sh)) * 0.065
	if size < 18 {
		return 18
	}
	if size > 46 {
		return 46
	}
	return size
}

func leaderboardHeaderFontSize(sw, sh int) float64 {
	size := float64(min(sw, sh)) * 0.026
	if size < 10 {
		return 10
	}
	if size > 18 {
		return 18
	}
	return size
}

func leaderboardRowFontSize(sw, sh int) float64 {
	size := float64(min(sw, sh)) * 0.03
	if size < 11 {
		return 11
	}
	if size > 20 {
		return 20
	}
	return size
}

func leaderboardHintFontSize(sw, sh int) float64 {
	size := float64(min(sw, sh)) * 0.026
	if size < 10 {
		return 10
	}
	if size > 18 {
		return 18
	}
	return size
}

func leaderboardRowHeight(sw, sh int) int {
	h := int(float64(sh) * 0.055)
	if h < 22 {
		return 22
	}
	if h > 38 {
		return 38
	}
	return h
}
