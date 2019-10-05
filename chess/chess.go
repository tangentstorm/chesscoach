package chess

import "fmt"

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

// Piece reperesents one of the 12 chess pieces, or NO for no piece.
type Piece = byte

const (
	// NO : marker for empty squares
	NO Piece = iota
	// WP : White Pawn
	WP
	// WR : White Rook
	WR
	// WN : White Knight
	WN
	// WB : White Bishop
	WB
	// WQ : WHite Queen
	WQ
	// WK : White King
	WK
	// BP : Black Pawn
	BP
	// BR : Black rook
	BR
	// BN : Black Knight
	BN
	// BB : Black Bishop
	BB
	// BQ : black queen
	BQ
	// BK : black knight
	BK
)

// A Board is a 2d array of pieces.
type Board = [8][8]Piece

// StartPos returns a Board in the usual starting position.
func StartPos() Board {
	return Board{
		{BR, BN, BB, BQ, BK, BB, BN, BR},
		{BP, BP, BP, BP, BP, BP, BP, BP},
		{NO, NO, NO, NO, NO, NO, NO, NO},
		{NO, NO, NO, NO, NO, NO, NO, NO},
		{NO, NO, NO, NO, NO, NO, NO, NO},
		{NO, NO, NO, NO, NO, NO, NO, NO},
		{WP, WP, WP, WP, WP, WP, WP, WP},
		{WR, WN, WB, WQ, WK, WB, WN, WR},
	}
}
