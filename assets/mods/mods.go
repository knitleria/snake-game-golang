package mods

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
	current = m
}

func Next() Mod {
	for i, m := range order {
		if m == current {
			return order[(i+1)%len(order)]
		}
	}
	return order[0]
}
