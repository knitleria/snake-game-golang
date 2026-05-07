package menu

import (
	"image/color"
	"snake_golang/assets"

	"github.com/hajimehoshi/ebiten/v2"
	etxt "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func Draw(dst *ebiten.Image, face *etxt.GoTextFaceSource, m *Menu) {
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

	if m == nil {
		return
	}
	m.layoutButtons()
	sw, sh := m.host.ScreenSize()
	if face == nil {
		return
	}
	drawMenuTitle(dst, face, sw, sh)
	drawMenuButtons(dst, face, m)
	drawMenuLangFlag(dst, m)
	drawMenuIcons(dst, m)
}

func titleFontSize(sw, sh int) float64 {
	m := sw
	if sh < m {
		m = sh
	}
	s := float64(m) * 0.075
	if s < 12 {
		s = 12
	}
	if s > 72 {
		s = 72
	}
	return s * 1.2
}

func drawMenuTitle(dst *ebiten.Image, face *etxt.GoTextFaceSource, sw, sh int) {
	size := titleFontSize(sw, sh)
	maxW := float64(sw) * 0.9
	size = fitFontSize(face, "Snake", size, maxW, 8)
	fo := &etxt.GoTextFace{Source: face, Size: size}
	op := &etxt.DrawOptions{}
	op.GeoM.Translate(float64(sw)/2, float64(sh)*0.2)
	op.PrimaryAlign = etxt.AlignCenter
	op.SecondaryAlign = etxt.AlignCenter
	op.ColorScale.ScaleWithColor(color.RGBA{R: 30, G: 55, B: 35, A: 255})
	etxt.Draw(dst, "Snake", fo, op)
}

func drawMenuIcons(dst *ebiten.Image, m *Menu) {
	for i, ic := range m.Icons {
		if ic.Image == nil {
			continue
		}
		r := ic.Rect
		iw := float64(ic.Image.Bounds().Dx())
		ih := float64(ic.Image.Bounds().Dy())
		if iw == 0 || ih == 0 {
			continue
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(r.Dx())/iw, float64(r.Dy())/ih)
		op.GeoM.Translate(float64(r.Min.X), float64(r.Min.Y))

		if i == m.IconFocus {
			op.ColorScale.Scale(1, 1, 1, 1)
		} else {
			op.ColorScale.Scale(0.85, 0.85, 0.85, 0.9)
		}
		dst.DrawImage(ic.Image, op)
	}
}

func drawMenuButtons(dst *ebiten.Image, face *etxt.GoTextFaceSource, m *Menu) {
	for i, b := range m.Buttons {
		rect := b.Rect
		x, y := float32(rect.Min.X), float32(rect.Min.Y)
		w, h := float32(rect.Dx()), float32(rect.Dy())

		fill := color.RGBA{R: 40, G: 40, B: 60, A: 220}
		if i == m.Focus {
			fill = color.RGBA{R: 80, G: 120, B: 200, A: 230}
		}
		vector.FillRect(dst, x, y, w, h, fill, true)
		vector.StrokeRect(dst, x, y, w, h, 2, color.RGBA{R: 240, G: 240, B: 240, A: 255}, true)

		label := b.Label()
		size := float64(rect.Dy()) * 0.45
		if size < 12 {
			size = 12
		}
		// Оставляем ~8% ширины кнопки на горизонтальные отступы.
		maxW := float64(rect.Dx()) * 0.92
		size = fitFontSize(face, label, size, maxW, 8)

		fo := &etxt.GoTextFace{Source: face, Size: size}
		op := &etxt.DrawOptions{}
		op.GeoM.Translate(
			float64(rect.Min.X)+float64(rect.Dx())/2,
			float64(rect.Min.Y)+float64(rect.Dy())/2,
		)
		op.PrimaryAlign = etxt.AlignCenter
		op.SecondaryAlign = etxt.AlignCenter
		op.ColorScale.ScaleWithColor(color.White)
		etxt.Draw(dst, label, fo, op)
	}
}

func drawMenuLangFlag(dst *ebiten.Image, m *Menu) {
	img := m.currentLangFlag()
	if img == nil {
		return
	}
	r := m.LangFlagRect
	iw := float64(img.Bounds().Dx())
	ih := float64(img.Bounds().Dy())
	if iw == 0 || ih == 0 {
		return
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(r.Dx())/iw, float64(r.Dy())/ih)
	op.GeoM.Translate(float64(r.Min.X), float64(r.Min.Y))

	if m.LangFlagHover {
		op.ColorScale.Scale(1, 1, 1, 1)
	} else {
		op.ColorScale.Scale(0.85, 0.85, 0.85, 0.9)
	}
	dst.DrawImage(img, op)
}

func fitFontSize(face *etxt.GoTextFaceSource, s string, size, maxW, minSize float64) float64 {
	if face == nil || s == "" || maxW <= 0 {
		return size
	}
	fo := &etxt.GoTextFace{Source: face, Size: size}
	w, _ := etxt.Measure(s, fo, 0)
	if w <= maxW {
		return size
	}
	scaled := size * maxW / w
	if scaled < minSize {
		scaled = minSize
	}
	return scaled
}
