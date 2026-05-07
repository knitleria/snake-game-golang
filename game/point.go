package game

type Point struct {
	X int
	Y int
}

var (
	Up    = Point{X: 0, Y: -1}
	Down  = Point{X: 0, Y: 1}
	Left  = Point{X: -1, Y: 0}
	Right = Point{X: 1, Y: 0}
)
