package skins

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
	current = m
}

func Next() Skin {
	for i, m := range order {
		if m == current {
			return order[(i+1)%len(order)]
		}
	}
	return order[0]
}
