package skins

import "snake_golang/game/unlock"

type Skin int

const (
	Normal Skin = iota
	Halfup
	Rantlol
)

func (m Skin) ID() string {
	switch m {
	case Halfup:
		return "halfup"
	case Rantlol:
		return "rantlol"
	default:
		return "normal"
	}
}

func (m Skin) LabelKey() string {
	switch m {
	case Halfup:
		return "mode.halfup"
	case Rantlol:
		return "mode.rantlol"
	default:
		return "mode.normal"
	}
}

var order = []Skin{Normal, Halfup, Rantlol}

var current = Normal

func Current() Skin {
	return current
}
func Set(m Skin) {
	if Locked(m) {
		return
	}
	current = m
}

func Next() Skin {
	n := len(order)
	idx := 0
	for i, m := range order {
		if m == current {
			idx = i
			break
		}
	}

	for step := 1; step <= n; step++ {
		cand := order[(idx+step)%n]
		if !Locked(cand) {
			return cand
		}
	}
	return current
}

func Locked(m Skin) bool {
	switch m {
	case Halfup, Rantlol:
		return !unlock.IsUnlocked()
	default:
		return false
	}
}
