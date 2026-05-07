package assets

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	_ "golang.org/x/image/webp"
)

//go:embed skins
var skinsFS embed.FS

//go:embed github.png soundcloud.png flags
var iconsFS embed.FS

var (
	themes         = map[string]*Theme{}
	active         *Theme
	IconGitHub     *ebiten.Image
	IconSoundCloud *ebiten.Image
	IconFlagRU     *ebiten.Image
	IconFlagEN     *ebiten.Image
	IconFlagUA     *ebiten.Image
)

func loadIcons() error {
	decode := func(name string) (*ebiten.Image, error) {
		b, err := iconsFS.ReadFile(name)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", name, err)
		}
		img, _, err := image.Decode(bytes.NewReader(b))
		if err != nil {
			return nil, fmt.Errorf("decode %s: %w", name, err)
		}
		return ebiten.NewImageFromImage(img), nil
	}
	var err error
	if IconGitHub, err = decode("github.png"); err != nil {
		return err
	}
	if IconSoundCloud, err = decode("soundcloud.png"); err != nil {
		return err
	}
	if IconFlagRU, err = decode("flags/flag_ru.png"); err != nil {
		return err
	}
	if IconFlagEN, err = decode("flags/flag_en.png"); err != nil {
		return err
	}
	if IconFlagUA, err = decode("flags/flag_ua.png"); err != nil {
		return err
	}
	return nil
}

func LoadAll(skinIDs []string) error {
	if err := loadIcons(); err != nil {
		return fmt.Errorf("load icons: %w", err)
	}
	for _, id := range skinIDs {
		t, err := loadTheme(id)
		if err != nil {
			return fmt.Errorf("load skin %q: %w", id, err)
		}
		themes[id] = t
	}
	return nil
}

func loadTheme(id string) (*Theme, error) {
	read := func(name string) ([]byte, error) {
		return skinsFS.ReadFile("skins/" + id + "/" + name)
	}
	img := func(name string) (*ebiten.Image, error) {
		b, err := read(name)
		if err != nil {
			return nil, err
		}
		decoded, _, err := image.Decode(bytes.NewReader(b))
		if err != nil {
			return nil, err
		}
		return ebiten.NewImageFromImage(decoded), nil
	}

	t := &Theme{}
	var err error
	if t.Background, err = img("bg.png"); err != nil {
		return nil, err
	}
	if t.MenuBackground, err = img("bg_menu.png"); err != nil {
		return nil, err
	}
	if t.Apple, err = img("apple.png"); err != nil {
		return nil, err
	}
	if t.SnakeHead, err = img("snake_head.png"); err != nil {
		return nil, err
	}
	if t.SnakeBody, err = img("snake_body.png"); err != nil {
		return nil, err
	}
	if t.SnakeTail, err = img("snake_tail.png"); err != nil {
		return nil, err
	}
	if t.SnakeDownRight, err = img("snake_down_right.png"); err != nil {
		return nil, err
	}
	if t.SnakeUpRight, err = img("snake_up_right.png"); err != nil {
		return nil, err
	}
	if t.SnakeLeftDown, err = img("snake_left_down.png"); err != nil {
		return nil, err
	}
	if t.SnakeLeftUp, err = img("snake_left_up.png"); err != nil {
		return nil, err
	}

	if t.GameMusicOGG, err = read("music_game.ogg"); err != nil {
		return nil, err
	}
	if t.MenuMusicOGG, err = read("music_menu.ogg"); err != nil {
		return nil, err
	}
	if t.EatSoundOGG, err = read("sfx_eat.ogg"); err != nil {
		return nil, err
	}

	return t, nil
}

func Current() *Theme {
	return active
}

func SetActive(skinID string) {
	if t, ok := themes[skinID]; ok {
		active = t
	}
}
