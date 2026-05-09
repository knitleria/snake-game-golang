package game

import (
	"errors"
	"image/color"
	"snake_golang/assets"
	"snake_golang/game/i18n"
	"snake_golang/game/profile"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	etxt "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (s *Screen) BeginNameInput() {
	if s.World == nil {
		return
	}
	s.nameDraft = []rune(s.PlayerName)
	s.nameError = ""
	s.World.State = StateNameInput
}

func (s *Screen) UpdateNameInput() error {
	for _, r := range ebiten.AppendInputChars(nil) {
		if r == '\n' || r == '\r' {
			continue
		}
		if len(s.nameDraft) < profile.MaxPlayerNameRunes {
			s.nameDraft = append(s.nameDraft, r)
			s.nameError = ""
		} else {
			s.nameError = i18n.T("name.error_long")
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(s.nameDraft) > 0 {
		s.nameDraft = s.nameDraft[:len(s.nameDraft)-1]
		s.nameError = ""
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if s.saveNameDraft() {
			s.World.State = StateMenu
		}
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) && s.PlayerName != "" {
		s.nameDraft = []rune(s.PlayerName)
		s.nameError = ""
		s.World.State = StateMenu
	}
	return nil
}

func (s *Screen) saveNameDraft() bool {
	name, err := profile.NormalizePlayerName(string(s.nameDraft))
	if err != nil {
		s.nameError = nameErrorHelper(err)
		return false
	}
	if err := profile.SaveConfig(profile.Config{PlayerName: name}); err != nil {
		s.nameError = nameErrorHelper(err)
		return false
	}
	s.PlayerName = name
	s.nameDraft = []rune(name)
	s.nameError = ""
	return true
}

func nameErrorHelper(err error) string {
	switch {
	case errors.Is(err, profile.ErrEmptyName):
		return i18n.T("name.error_empty")
	case errors.Is(err, profile.ErrNameTooLong):
		return i18n.T("name.error_long")
	case errors.Is(err, profile.ErrBadName):
		return i18n.T("name.error_invalid")
	default:
		return i18n.T("name.error_save")
	}
}

func DrawNameInput(dst *ebiten.Image, face *etxt.GoTextFaceSource, name string, errText string, canCancel bool) {
	drawNameInputBackground(dst)

	if face == nil {
		return
	}

	sw, sh := dst.Bounds().Dx(), dst.Bounds().Dy()

	title := i18n.T("name.title")
	input := name
	if input == "" {
		input = i18n.T("name.placeholder")
	}
	input += "_"

	hint := i18n.T("name.hint_required")
	if canCancel {
		hint = i18n.T("name.hint_optional")
	}

	drawNameText(dst, face, title, float64(sw)/2, float64(sh)*0.30, nameTitleFontSize(sw, sh), color.White)
	drawNameInputBox(dst, face, input, sw, sh)
	drawNameText(dst, face, hint, float64(sw)/2, float64(sh)*0.62, nameHintFontSize(sw, sh), color.RGBA{R: 220, G: 225, B: 230, A: 255})

	if errText != "" {
		drawNameText(dst, face, errText, float64(sw)/2, float64(sh)*0.70, nameHintFontSize(sw, sh), color.RGBA{R: 255, G: 130, B: 130, A: 255})
	}
}

func drawNameInputBackground(dst *ebiten.Image) {
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
	vector.FillRect(dst, 0, 0, float32(dw), float32(dh), color.RGBA{R: 20, G: 22, B: 28, A: 150}, true)
}

func drawNameInputBox(dst *ebiten.Image, face *etxt.GoTextFaceSource, text string, sw, sh int) {
	boxW := int(float64(sw) * 0.52)
	if boxW < 260 {
		boxW = 260
	}
	if maxW := sw - 40; boxW > maxW {
		boxW = maxW
	}

	boxH := int(float64(sh) * 0.11)
	if boxH < 54 {
		boxH = 54
	}
	if boxH > 86 {
		boxH = 86
	}

	x := (sw - boxW) / 2
	y := int(float64(sh) * 0.42)

	vector.FillRect(dst, float32(x), float32(y), float32(boxW), float32(boxH),
		color.RGBA{R: 35, G: 38, B: 48, A: 230}, true)
	vector.StrokeRect(dst, float32(x), float32(y), float32(boxW), float32(boxH), 2,
		color.RGBA{R: 240, G: 240, B: 245, A: 255}, true)

	size := float64(boxH) * 0.42
	size = FitFontSize(face, text, size, float64(boxW)*0.88, 10)
	drawNameText(dst, face, text, float64(sw)/2, float64(y)+float64(boxH)/2, size, color.White)
}

func drawNameText(dst *ebiten.Image, face *etxt.GoTextFaceSource, text string, x, y, size float64, clr color.Color) {
	fo := &etxt.GoTextFace{Source: face, Size: size}
	op := &etxt.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.PrimaryAlign = etxt.AlignCenter
	op.SecondaryAlign = etxt.AlignCenter
	op.ColorScale.ScaleWithColor(clr)
	etxt.Draw(dst, text, fo, op)
}

func nameTitleFontSize(sw, sh int) float64 {
	m := min(sw, sh)
	size := float64(m) * 0.07
	if size < 18 {
		return 18
	}
	if size > 48 {
		return 48
	}
	return size
}

func nameHintFontSize(sw, sh int) float64 {
	m := min(sw, sh)
	size := float64(m) * 0.032
	if size < 11 {
		return 11
	}
	if size > 22 {
		return 22
	}
	return size
}

func DrawMenuPlayerName(dst *ebiten.Image, face *etxt.GoTextFaceSource, playerName string) {
	if face == nil {
		return
	}

	sw, sh := dst.Bounds().Dx(), dst.Bounds().Dy()
	name := playerName
	if name == "" {
		name = i18n.T("name.not_set")
	}

	label := i18n.T("menu.player") + ": " + name
	hint := i18n.T("name.edit_hint")

	padding := float64(HudPadding)
	size := menuPlayerNameFontSize(sw, sh)

	fo := &etxt.GoTextFace{Source: face, Size: size}
	op := &etxt.DrawOptions{}
	op.GeoM.Translate(padding, padding)
	op.PrimaryAlign = etxt.AlignStart
	op.SecondaryAlign = etxt.AlignStart
	op.ColorScale.ScaleWithColor(color.RGBA{R: 245, G: 245, B: 245, A: 245})
	etxt.Draw(dst, label, fo, op)

	hintSize := size * 0.68
	hintFo := &etxt.GoTextFace{Source: face, Size: hintSize}
	hintOp := &etxt.DrawOptions{}
	hintOp.GeoM.Translate(padding, padding+size*1.25)
	hintOp.PrimaryAlign = etxt.AlignStart
	hintOp.SecondaryAlign = etxt.AlignStart
	hintOp.ColorScale.ScaleWithColor(color.RGBA{R: 220, G: 220, B: 220, A: 210})
	etxt.Draw(dst, hint, hintFo, hintOp)
}

func menuPlayerNameFontSize(sw, sh int) float64 {
	m := min(sw, sh)
	size := float64(m) * 0.026
	if size < 11 {
		return 11
	}
	if size > 20 {
		return 20
	}
	return size
}
