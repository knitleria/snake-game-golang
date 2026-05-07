package assets

import "github.com/hajimehoshi/ebiten/v2"

const SampleRate = 44100

type Theme struct {
	Background     *ebiten.Image
	MenuBackground *ebiten.Image
	Apple          *ebiten.Image
	SnakeHead      *ebiten.Image
	SnakeBody      *ebiten.Image
	SnakeTail      *ebiten.Image
	SnakeDownRight *ebiten.Image
	SnakeUpRight   *ebiten.Image
	SnakeLeftDown  *ebiten.Image
	SnakeLeftUp    *ebiten.Image

	GameMusicOGG []byte
	MenuMusicOGG []byte
	EatSoundOGG  []byte
}
