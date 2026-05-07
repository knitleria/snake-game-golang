package game

import (
	"math"
	"snake_golang/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

// Базовые ориентации спрайтов.
// Голова — язык влево, база = Left (π).
// Хвост — остриё направлено вправо, тело слева, база = Left (π).
// Тело — горизонтальное, база = Right (0).
const (
	baseAngleHead = math.Pi
	baseAngleTail = math.Pi
	baseAngleBody = 0
)

func drawSnake(dst *ebiten.Image, w *World, cw, ch float32) {
	n := len(w.Snake)
	for i, p := range w.Snake {
		sprite, angle := snakeSegmentSprite(w, i, n)
		if sprite == nil {
			continue
		}
		drawCellSprite(dst, sprite, p, angle, cw, ch)
	}
}

func snakeSegmentSprite(w *World, i, n int) (*ebiten.Image, float64) {
	active := assets.Current()
	if active == nil {
		return nil, 0
	}
	switch {
	case i == 0:
		dir := w.Direction
		if n > 1 {
			dir = segmentDir(w, w.Snake[0], w.Snake[1])
		}
		return active.SnakeHead, dirAngle(dir) - baseAngleHead
	case i == n-1:
		dir := segmentDir(w, w.Snake[i-1], w.Snake[i])
		return active.SnakeTail, dirAngle(dir) - baseAngleTail
	default:
		inDir := segmentDir(w, w.Snake[i], w.Snake[i+1])
		outDir := segmentDir(w, w.Snake[i-1], w.Snake[i])
		if inDir == outDir {
			return active.SnakeBody, dirAngle(inDir) - baseAngleBody
		}
		return turnSprite(inDir, outDir), 0
	}
}

// turnSprite выбирает спрайт угла по направлению движения до/после поворота.
// Имена картинок обозначают стороны, с которых у угла открытые концы;
// противоположные направления движения дают тот же набор открытий и используют тот же спрайт.
// Открытия = { -inDir, outDir }.
func turnSprite(inDir, outDir Point) *ebiten.Image {
	active := assets.Current()
	if active == nil {
		return nil
	}
	switch {
	// ┌ : открытия снизу и справа (движение Up→Right или Left→Down)
	case inDir == Up && outDir == Right,
		inDir == Left && outDir == Down:
		return active.SnakeDownRight
	// └ : открытия сверху и справа (Down→Right или Left→Up)
	case inDir == Down && outDir == Right,
		inDir == Left && outDir == Up:
		return active.SnakeUpRight
	// ┐ : открытия слева и снизу (Right→Down или Up→Left)
	case inDir == Right && outDir == Down,
		inDir == Up && outDir == Left:
		return active.SnakeLeftDown
	// ┘ : открытия слева и сверху (Right→Up или Down→Left)
	case inDir == Right && outDir == Up,
		inDir == Down && outDir == Left:
		return active.SnakeLeftUp
	}
	return active.SnakeBody
}

func drawCellSprite(dst *ebiten.Image, sprite *ebiten.Image, p Point, angle float64, cw, ch float32) {
	b := sprite.Bounds()
	sw, sh := b.Dx(), b.Dy()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(sw)/2, -float64(sh)/2)
	op.GeoM.Scale(float64(cw)/float64(sw), float64(ch)/float64(sh))
	op.GeoM.Rotate(angle)
	op.GeoM.Translate(
		float64(p.X)*float64(cw)+float64(cw)/2,
		float64(p.Y)*float64(ch)+float64(ch)/2,
	)
	dst.DrawImage(sprite, op)
}

func subPoint(a, b Point) Point { return Point{X: a.X - b.X, Y: a.Y - b.Y} }

func segmentDir(w *World, a, b Point) Point {
	d := subPoint(a, b)
	if w == nil || !w.WrapEdges {
		return d
	}

	if d.X == -(GridWidth - 1) {
		return Right
	}
	if d.X == GridWidth-1 {
		return Left
	}
	if d.Y == -(GridHeight - 1) {
		return Down
	}
	if d.Y == GridHeight-1 {
		return Up
	}
	return d
}

// dirAngle — угол поворота спрайта в ebiten-конвенции (положительный = по часовой).
func dirAngle(d Point) float64 {
	switch d {
	case Right:
		return 0
	case Down:
		return math.Pi / 2
	case Left:
		return math.Pi
	case Up:
		return -math.Pi / 2
	}
	return 0
}
