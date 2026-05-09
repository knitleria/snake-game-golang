package main

import (
	"bytes"
	"log"

	"snake_golang/assets"
	"snake_golang/assets/skins"
	game "snake_golang/game"
	"snake_golang/game/profile"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	etxt "github.com/hajimehoshi/ebiten/v2/text/v2"
)

func main() {
	faceSource, err := etxt.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	w := game.NewWorld()
	config, err := profile.LoadConfig()
	if err != nil {
		log.Printf("load config: %v", err)
	}
	if config.PlayerName == "" {
		w.State = game.StateNameInput
	}

	ebiten.SetFullscreen(false)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)

	baseW := game.GridWidth * game.DefaultCellSize
	baseH := game.GridHeight * game.DefaultCellSize
	winW, winH := baseW, baseH
	posX, posY := -1, -1
	if m := ebiten.Monitor(); m != nil {
		if mw, mh := m.Size(); mw > 0 && mh > 0 {
			scale := float64(mw) * game.WindowMonitorFraction / float64(baseW)
			if sh := float64(mh) * game.WindowMonitorFraction / float64(baseH); sh < scale {
				scale = sh
			}
			if scale > 1 {
				winW = int(float64(baseW) * scale)
				winH = int(float64(baseH) * scale)
			}
			posX = (mw - winW) / 2
			posY = (mh - winH) / 2
		}
	}
	ebiten.SetWindowSize(winW, winH)
	if posX >= 0 && posY >= 0 {
		ebiten.SetWindowPosition(posX, posY)
	}

	ebiten.SetWindowTitle("Snake Game")

	if err := assets.LoadAll([]string{
		skins.Normal.ID(),
		skins.Halfup.ID(),
		skins.Rantlol.ID(),
	}); err != nil {
		log.Fatal(err)
	}
	assets.SetActive(skins.Current().ID())

	audioContext := audio.NewContext(assets.SampleRate)
	gameAudio := game.NewAudio(audioContext)
	if err := gameAudio.Reload(assets.Current()); err != nil {
		log.Fatal(err)
	}

	screen := &game.Screen{
		World:      w,
		FaceSource: faceSource,
		Audio:      gameAudio,
		PlayerName: config.PlayerName,
	}
	if err := ebiten.RunGame(screen); err != nil {
		log.Fatal(err)
	}
}
