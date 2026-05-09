package game

import (
	"snake_golang/assets/mods"
	"time"
)

type State int

const (
	StateMenu State = iota
	StateWaiting
	StatePlaying
	StatePaused
	StateGameOver
	StateNameInput
)

func (s State) String() string {
	switch s {
	case StateMenu:
		return "Menu"
	case StateWaiting:
		return "Waiting"
	case StatePlaying:
		return "Playing"
	case StatePaused:
		return "Paused"
	case StateGameOver:
		return "GameOver"
	case StateNameInput:
		return "NameInput"
	default:
		return "State(unknown)"
	}
}

type World struct {
	State         State
	Snake         []Point
	Direction     Point
	NextDirection Point
	LastMove      time.Time
	Apple         Point
	Score         int
	WrapEdges     bool
}

func (w *World) IsPlaying() bool {
	return w.State == StatePlaying
}

func (w *World) moveInterval() time.Duration {
	n := len(w.Snake)
	step := BaseInterval - time.Duration(n-1)*PerSegment
	if step < MinInterval {
		return MinInterval
	}
	return step
}

func NewWorld() *World {
	w := &World{
		State:         StateMenu,
		Snake:         []Point{{X: GridWidth / 2, Y: GridHeight / 2}},
		Direction:     Point{X: 0, Y: 1},
		NextDirection: Point{X: 0, Y: 1},
		LastMove:      time.Now(),
		Score:         0,
	}
	GenerateApple(w)
	return w
}

func newWorldPlaying() *World {
	w := NewWorld()
	w.WrapEdges = mods.Current() == mods.Defaltyk
	w.State = StatePlaying
	w.LastMove = time.Now()
	return w
}
