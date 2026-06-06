package mods

import "snake_golang/game/unlock"

type Mod int

const (
	Normal Mod = iota
	Defaltyk
)

func (m Mod) ID() string {
	switch m {
	case Defaltyk:
		return "defaltyk"
	default:
		return "normal"
	}
}

func (m Mod) LabelKey() string {
	switch m {
	case Defaltyk:
		return "mode.defaltyk"
	default:
		return "mode.normal"
	}
}

var order = []Mod{Normal, Defaltyk}

var current = Normal

func Current() Mod {
	return current
}
func Set(m Mod) {
	if Locked(m) {
		return
	}
	current = m
}

func Next() Mod {
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

func Locked(m Mod) bool {
	switch m {
	case Defaltyk:
		return !unlock.IsUnlocked()
	default:
		return false
	}
}
