package game

import (
	"math/rand"
)

func pointOnSnake(p Point, snake []Point) bool {
	for _, s := range snake {
		if s == p {
			return true
		}
	}
	return false
}

func Step(w *World) {
	w.Direction = w.NextDirection
	head := w.Snake[0]
	newHead := Point{X: head.X + w.Direction.X, Y: head.Y + w.Direction.Y}

	if w.WrapEdges {
		newHead = WrapPoint(newHead)
	}

	if Collision(newHead, w.Snake) {
		w.State = StateGameOver
		return
	}
	if newHead == w.Apple {
		w.Snake = append([]Point{newHead}, w.Snake...)
		w.Score += ScorePerApple
		GenerateApple(w)
		return
	}
	w.Snake = append([]Point{newHead}, w.Snake[:len(w.Snake)-1]...)
}

func GenerateApple(w *World) {
	if len(w.Snake) >= GridWidth*GridHeight {
		return
	}
	for {
		p := Point{X: rand.Intn(GridWidth), Y: rand.Intn(GridHeight)}
		if !pointOnSnake(p, w.Snake) {
			w.Apple = p
			return
		}
	}
}

func WrapPoint(p Point) Point {
	if p.X < 0 {
		p.X = GridWidth - 1
	} else if p.X >= GridWidth {
		p.X = 0
	}
	if p.Y < 0 {
		p.Y = GridHeight - 1
	} else if p.Y >= GridHeight {
		p.Y = 0
	}
	return p
}

func IsOppositeDirection(a, b Point) bool {
	return a.X == -b.X && a.Y == -b.Y
}
