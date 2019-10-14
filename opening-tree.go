package main

import (
	"fmt"
	"log"

	"github.com/notnil/chess"

	"github.com/notnil/opening"
	uci "gopkg.in/freeeve/uci.v1"
)

// MoveScore is the engine's score for a possible next move.
type MoveScore struct {
	// Move is the chess move.
	Move *chess.Move
	// Note is a string containing the algebraic notation for the move.
	Note string
	// Name is the name of the position it creates, if any.
	Name string
	// Score is the score the engine assigned to the move.
	Score int
}

// Coach combines a game and engine
type Coach struct {
	game *chess.Game
	eng  *uci.Engine
}

// NewCoach creates a new coach
func NewCoach(enginePath string) (c *Coach, e error) {
	c = &Coach{
		game: chess.NewGame(chess.UseNotation(chess.AlgebraicNotation{})),
	}
	c.eng, e = uci.NewEngine(enginePath)
	if e != nil {
		return nil, e
	}
	c.eng.SetOptions(uci.Options{
		Hash:    128,
		Ponder:  false,
		OwnBook: true,
		MultiPV: 32,
	})
	return
}

// Clone creates a copy of the coach with its own game history.
func (coach *Coach) Clone() *Coach {
	return &Coach{
		game: coach.game.Clone(),
		eng:  coach.eng,
	}
}

// PGNHistory returns the moves in PGN notation.
func (coach *Coach) PGNHistory() string {
	s := ""
	// !! this was basically copied from chess/pgn:encodePGN.
	// TODO: make this accessible as a method in chess itself.
	g := coach.game
	for i, move := range g.Moves() {
		pos := g.Positions()[i]
		txt := (chess.AlgebraicNotation{}).Encode(pos, move)
		if i%2 == 0 {
			s += fmt.Sprintf("%d.%s", (i/2)+1, txt)
		} else {
			s += fmt.Sprintf(" %s ", txt)
		}
	}
	return s
}

func indent(i int) (result string) {
	result = ""
	for j := 0; j < i; j++ {
		result += " "
	}
	return
}

// BestMoves returns the list of n (or fewer) best moves. If n=0, ranks all moves.
func (coach *Coach) BestMoves(n int) (result []MoveScore) {
	coach.eng.SetFEN(coach.game.Position().String())
	opts := uci.HighestDepthOnly | uci.IncludeUpperbounds | uci.IncludeLowerbounds
	search, _ := coach.eng.GoDepth(16, opts)
	for i, r := range search.Results {
		if i < n || n == 0 {
			best := r.BestMoves[0]
			move, err := (chess.LongAlgebraicNotation{}).Decode(coach.game.Position(), best)
			quitOn(err)
			name := ""
			if opn := opening.Find(coach.game.Moves()); opn != nil {
				name = fmt.Sprintf("%s %s", opn.Code(), opn.Title())
			}
			result = append(result, MoveScore{
				Move:  move,
				Note:  (chess.AlgebraicNotation{}).Encode(coach.game.Position(), move),
				Name:  name,
				Score: r.Score,
			})
		}
	}
	return result
}

// walk the game tree
func (coach *Coach) walk(depth int, mv *chess.Move) {
	c := coach
	if depth > 0 {
		c = c.Clone()
		err := c.game.Move(mv)
		quitOn(err)
	}
	if depth < 5 {
		for _, best := range c.BestMoves(2) {
			fmt.Printf("%s [%s] -> %s (score: %d)\n", c.PGNHistory(), best.Name, best.Note, best.Score)
			c.walk(depth+1, best.Move)
		}
	}
}

func quitOn(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	coach, err := NewCoach("./stockfish_10_x64")
	quitOn(err)
	coach.walk(0, nil)
}
