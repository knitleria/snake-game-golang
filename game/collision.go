package game

import "slices"

func Collision(head Point, snake []Point) bool {
	if head.X < 0 || head.X >= GridWidth || head.Y < 0 || head.Y >= GridHeight {
		return true
	}
	if slices.Contains(snake, head) {
		return true
	}
	return false
}
