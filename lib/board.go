package chess

import (
	"fmt"
)

// A Square is a point on the chessboard
type Square struct{ x, y int }

// Nowhere represents a square that is not on the board.
var Nowhere = Square{-1, -1}

// SquareAt converts coordinates to a Square
func SquareAt(x, y int) Square {
	if x >= 0 && y >= 0 && x <= 7 && y <= 7 {
		return Square{x, y}
	}
	return Nowhere
}

var files = []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'}

// Name is the name of the square in algebraic notation (e.g, "e4").
func (s Square) Name() string {
	return fmt.Sprintf("%c%d", files[s.x], 8-s.y)
}
