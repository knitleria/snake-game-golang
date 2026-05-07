package menu

import (
	"image"
	"snake_golang/assets"
	"snake_golang/assets/mods"
	"snake_golang/assets/skins"
	"snake_golang/game/i18n"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Host interface {
	StartGame()
	ScreenSize() (int, int)
	SwitchSkin()
	SwitchMode()
}

type IconButton struct {
	Rect  image.Rectangle
	Image *ebiten.Image
	URL   string
}

type LangButton struct {
	Image *ebiten.Image
	Lang  i18n.Lang
}

type Button struct {
	Rect    image.Rectangle
	Label   func() string
	OnClick func()
}

type SkinButton struct {
	Rect  image.Rectangle
	Image *ebiten.Image
	Skin  skins.Skin
}

type Menu struct {
	Buttons       []Button
	SkinsButtons  []SkinButton
	Icons         []IconButton
	Focus         int
	IconFocus     int
	LangButtons   []LangButton
	LangFlagRect  image.Rectangle
	LangFlagHover bool
	host          Host
}

func NewMenu(h Host) *Menu {
	return &Menu{
		host: h,
		Buttons: []Button{
			{
				Label:   func() string { return i18n.T("menu.play") },
				OnClick: h.StartGame,
			},
			{
				Label: func() string {
					return i18n.T("menu.mode") + ": " + i18n.T(modeKey(mods.Current()))
				},
				OnClick: h.SwitchMode,
			},
			{
				Label: func() string {
					return i18n.T("menu.skins") + ": " + i18n.T(skinsKey(skins.Current()))
				},
				OnClick: h.SwitchSkin,
			},
		},
		Icons: []IconButton{
			{Image: assets.IconGitHub, URL: "https://github.com/yourusername/yourrepo"},
			{Image: assets.IconSoundCloud, URL: "https://soundcloud.com/halfup"},
		},
		LangButtons: []LangButton{
			{Image: assets.IconFlagRU, Lang: i18n.LangRU},
			{Image: assets.IconFlagEN, Lang: i18n.LangEN},
			{Image: assets.IconFlagUA, Lang: i18n.LangUA},
		},
	}
}

func (m *Menu) currentLangFlag() *ebiten.Image {
	cur := i18n.CurrentLang()
	for _, b := range m.LangButtons {
		if b.Lang == cur {
			return b.Image
		}
	}
	return nil
}

func (m *Menu) cycleLang() {
	if len(m.LangButtons) == 0 {
		return
	}
	cur := i18n.CurrentLang()
	for i, b := range m.LangButtons {
		if b.Lang == cur {
			next := m.LangButtons[(i+1)%len(m.LangButtons)].Lang
			i18n.SetLang(next)
			return
		}
	}
	i18n.SetLang(m.LangButtons[0].Lang)
}

func langKey(l i18n.Lang) string {
	switch l {
	case i18n.LangRU:
		return "lang.ru"
	case i18n.LangEN:
		return "lang.en"
	case i18n.LangUA:
		return "lang.ua"
	}
	return "lang.en"
}

func modeKey(m mods.Mod) string {
	switch m {
	case mods.Normal:
		return "mode.normal"
	case mods.Defaltyk:
		return "mode.defaltyk"
	}
	return "mode.normal"
}

func skinsKey(s skins.Skin) string {
	switch s {
	case skins.Normal:
		return "skins.normal"
	case skins.Halfup:
		return "skins.halfup"
	case skins.Rantlol:
		return "skins.rantlol"
	}
	return "skins.normal"
}
func (m *Menu) layoutIcons() {
	sw, sh := m.host.ScreenSize()

	// Размер иконки: ~6% от меньшей стороны экрана
	size := int(float64(min(sw, sh)) * 0.06)
	if size < 28 {
		size = 28
	}
	if size > 64 {
		size = 64
	}

	margin := size / 2
	gap := size / 3
	x := margin
	y := sh - margin - size

	for i := range m.Icons {
		m.Icons[i].Rect = image.Rect(x, y, x+size, y+size)
		x += size + gap
	}
}

func (m *Menu) layoutLangFlag() {
	sw, sh := m.host.ScreenSize()

	// Размер флага: ~6% от меньшей стороны экрана
	size := int(float64(min(sw, sh)) * 0.06)
	if size < 28 {
		size = 28
	}
	if size > 64 {
		size = 64
	}

	margin := size / 2
	x := sw - margin - size
	y := margin
	m.LangFlagRect = image.Rect(x, y, x+size, y+size)
}

func (m *Menu) layoutButtons() {
	n := len(m.Buttons)
	if n == 0 {
		return
	}
	sw, sh := m.host.ScreenSize()

	bw := int(float64(sw) * 0.4)
	if bw < 200 {
		bw = 200
	}
	if maxW := sw - 40; bw > maxW {
		bw = maxW
	}
	if bw < 1 {
		bw = 1
	}

	bh := int(float64(sh) * 0.1)
	if bh < 48 {
		bh = 48
	}

	gap := int(float64(bh) * 0.3)
	total := n*bh + (n-1)*gap
	startY := (sh-total)/2 + int(float64(sh)*0.08)
	x := (sw - bw) / 2

	for i := range m.Buttons {
		y := startY + i*(bh+gap)
		m.Buttons[i].Rect = image.Rect(x, y, x+bw, y+bh)
	}
}

func (m *Menu) Update() error {
	m.layoutButtons()
	n := len(m.Buttons)
	if n == 0 {
		return nil
	}
	if m.Focus < 0 || m.Focus >= n {
		m.Focus = 0
	}

	cx, cy := ebiten.CursorPosition()
	cursor := image.Point{X: cx, Y: cy}
	for i := range m.Buttons {
		if cursor.In(m.Buttons[i].Rect) {
			m.Focus = i
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				m.Buttons[i].OnClick()
				return nil
			}
			break
		}
	}
	m.layoutLangFlag()
	m.LangFlagHover = cursor.In(m.LangFlagRect)
	if m.LangFlagHover && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		m.cycleLang()
		return nil
	}

	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyUp), inpututil.IsKeyJustPressed(ebiten.KeyW):
		m.Focus = (m.Focus - 1 + n) % n
	case inpututil.IsKeyJustPressed(ebiten.KeyDown), inpututil.IsKeyJustPressed(ebiten.KeyS):
		m.Focus = (m.Focus + 1) % n
	case inpututil.IsKeyJustPressed(ebiten.KeyEnter), inpututil.IsKeyJustPressed(ebiten.KeySpace):
		m.Buttons[m.Focus].OnClick()
	}
	m.layoutIcons()
	m.IconFocus = -1
	for i := range m.Icons {
		if cursor.In(m.Icons[i].Rect) {
			m.IconFocus = i
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				openURL(m.Icons[i].URL)
				return nil
			}
			break
		}
	}
	return nil
}
