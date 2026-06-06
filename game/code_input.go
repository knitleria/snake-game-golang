package game

import (
	"image/color"
	"snake_golang/assets/mods"
	"snake_golang/assets/skins"
	"snake_golang/game/i18n"
	"snake_golang/game/unlock"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	etxt "github.com/hajimehoshi/ebiten/v2/text/v2"
)

const maxCodeRunes = 32

func (s *Screen) BeginCodeInput() {
	if s.World == nil {
		return
	}
	s.codeDraft = []rune{}
	s.codeError = ""
	s.World.State = StateCodeInput
}

func (s *Screen) UpdateCodeInput() error {

	for _, r := range ebiten.AppendInputChars(nil) {
		if r == '\n' || r == '\r' {
			continue
		}
		if len(s.codeDraft) < maxCodeRunes {
			s.codeDraft = append(s.codeDraft, r)
			s.codeError = ""
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(s.codeDraft) > 0 {
		s.codeDraft = s.codeDraft[:len(s.codeDraft)-1]
		s.codeError = ""
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if unlock.TryUnlock(string(s.codeDraft)) {
			s.codeDraft = nil
			s.codeError = ""
			s.World.State = StateMenu
		} else {
			s.codeError = i18n.T("code.error_wrong")
		}
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		s.codeDraft = nil
		s.codeError = ""
		s.World.State = StateMenu
	}
	return nil
}

func DrawCodeInput(dst *ebiten.Image, face *etxt.GoTextFaceSource, code string, errText string) {
	drawNameInputBackground(dst)

	if face == nil {
		return
	}

	sw, sh := dst.Bounds().Dx(), dst.Bounds().Dy()

	title := i18n.T("code.title")
	input := code
	if input == "" {
		input = i18n.T("code.placeholder")
	}
	input += "_"

	hint := i18n.T("code.hint")

	drawNameText(dst, face, title, float64(sw)/2, float64(sh)*0.30, nameTitleFontSize(sw, sh), color.White)
	drawNameInputBox(dst, face, input, sw, sh)
	drawNameText(dst, face, hint, float64(sw)/2, float64(sh)*0.62, nameHintFontSize(sw, sh), color.RGBA{R: 220, G: 225, B: 230, A: 255})

	if errText != "" {
		drawNameText(dst, face, errText, float64(sw)/2, float64(sh)*0.70, nameHintFontSize(sw, sh), color.RGBA{R: 255, G: 130, B: 130, A: 255})
	}
}

func DrawMenuCodeHint(dst *ebiten.Image, face *etxt.GoTextFaceSource) {
	if face == nil {
		return
	}
	if !skins.Locked(skins.Halfup) || !skins.Locked(skins.Rantlol) && !mods.Locked(mods.Defaltyk) {
		return
	}

	sw, sh := dst.Bounds().Dx(), dst.Bounds().Dy()

	hint := i18n.T("code.menu_hint")

	padding := float64(HudPadding)
	nameSize := menuPlayerNameFontSize(sw, sh)
	size := nameSize * 0.68

	y := padding + nameSize*1.25 + nameSize*1.25

	fo := &etxt.GoTextFace{Source: face, Size: size}
	op := &etxt.DrawOptions{}
	op.GeoM.Translate(padding, y)
	op.PrimaryAlign = etxt.AlignStart
	op.SecondaryAlign = etxt.AlignStart
	op.ColorScale.ScaleWithColor(color.RGBA{R: 220, G: 225, B: 230, A: 255})
	etxt.Draw(dst, hint, fo, op)

}
