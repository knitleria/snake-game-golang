package game

import (
	"fmt"
	"image/color"
	"snake_golang/game/i18n"

	"github.com/hajimehoshi/ebiten/v2"
	etxt "github.com/hajimehoshi/ebiten/v2/text/v2"
)

func HudFontSize() float64 {
	m := min(ScreenWidth, ScreenHeight)
	s := float64(m) * 0.05
	if s < HudFontMinSize {
		return HudFontMinSize
	}
	if s > HudFontMaxSize {
		return HudFontMaxSize
	}
	return s
}

func DrawScoreHud(dst *ebiten.Image, face *etxt.GoTextFaceSource, w *World) {
	if dst == nil || face == nil || w == nil {
		return
	}

	display := int(w.Score)
	text := fmt.Sprintf(i18n.T("game.score"), display)

	size := HudFontSize()
	maxW := float64(ScreenWidth - 2*HudPadding)
	size = FitFontSize(face, text, size, maxW, HudFontMinSize)
	fo := &etxt.GoTextFace{Source: face, Size: size}

	x, y := float64(HudPadding), float64(HudPadding)

	shadow := &etxt.DrawOptions{}
	shadow.GeoM.Translate(x+2, y+2)
	shadow.PrimaryAlign = etxt.AlignStart
	shadow.SecondaryAlign = etxt.AlignStart
	shadow.ColorScale.ScaleWithColor(color.RGBA{R: 0, G: 0, B: 0, A: 180})
	etxt.Draw(dst, text, fo, shadow)

	op := &etxt.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.PrimaryAlign = etxt.AlignStart
	op.SecondaryAlign = etxt.AlignStart
	op.ColorScale.ScaleWithColor(color.RGBA{R: 255, G: 230, B: 120, A: 255})
	etxt.Draw(dst, text, fo, op)
}
