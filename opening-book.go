package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/mattn/go-sqlite3"
	"github.com/notnil/chess"
	"github.com/notnil/opening"
	uci "gopkg.in/freeeve/uci.v1"
)

func quitOn(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type OpeningBook struct {
	db *sql.DB
}

func OpenOB(path string) (ob *OpeningBook, err error) {
	db, err := sql.Open("sqlite3", path)
	if err == nil {
		_, err = db.Exec(`
      create table if not exists position (
      id integer not null primary key autoincrement,
      board text );

      -- starting position
      insert or ignore into position (id, board) values
        (0, 'rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR');

      create table if not exists move (
        id integer primary key autoincrement,
        prev integer not null default 0,
        move text,
        score integer,
        unique(prev, move));`)
		if err == nil {
			ob = &OpeningBook { db }
		}
	}
	return ob, err
}

func (ob *OpeningBook) AddLine(prev int, move string) (id int) {
	s,e := ob.db.Prepare(`insert or ignore into move (prev, move) values (?, ?)`)
	quitOn(e)
	s.Exec(prev, move)
	q,e := ob.db.Prepare(`select id from move where prev=? and move=?`)
	quitOn(e)
	e = q.QueryRow(prev, move).Scan(&id);
	return
}

func (ob *OpeningBook) ScoreLine(prev int, move string, score int) (id int) {
	id = ob.AddLine(prev, move)
	s,e := ob.db.Prepare(`update move set score=? where id=?`)
	quitOn(e)
	s.Exec(score, id)
	quitOn(e)
	return
}

func (ob *OpeningBook) Line(n int) (game *chess.Game, err error) {
	stmt, err := ob.db.Prepare(`
	with recursive line(pre) as
    (values(?) union select prev from move, line where id=line.pre)
  select id, move from move where id in line`)
	quitOn(err)
	if err == nil {
		defer stmt.Close()
		game = chess.NewGame()
		rows, err := stmt.Query(n)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var id int
			var mv string
			err = rows.Scan(&id, &mv)
			if err != nil {
				return nil, err
			}
			game.MoveStr(mv)
		}
		if err = rows.Err(); err != nil {
			return nil, err
		}
	}
	return game, err
}

func (ob *OpeningBook) Close() {
	if ob != nil && ob.db != nil {
		ob.db.Close()
	}
}


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
	ob   *OpeningBook
}

// NewCoach creates a new coach
func NewCoach(enginePath string, bookPath string) (c *Coach, e error) {
	c = &Coach{
		game: chess.NewGame(chess.UseNotation(chess.AlgebraicNotation{})),
	}
	c.eng, e = uci.NewEngine(enginePath)
	if e == nil {
		c.eng.SetOptions(uci.Options{
			Hash:    128,
			Ponder:  false,
			OwnBook: true,
			MultiPV: 8,
		})
		c.ob, e = OpenOB(bookPath)
	}
	if e != nil {
		return nil, e
	}
	return
}

// Clone creates a copy of the coach with its own game history.
func (coach *Coach) Clone() *Coach {
	return &Coach{
		game: coach.game.Clone(),
		eng:  coach.eng,
		ob:   coach.ob,
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
func (coach *Coach) walk(prev int, depth int, mv *chess.Move) {
	c := coach
	if depth > 0 {
		c = c.Clone()
		err := c.game.Move(mv)
		quitOn(err)
	}
	if depth < 8 {
		for _, best := range c.BestMoves(4) {
			id := c.ob.ScoreLine(prev, best.Note, best.Score)
			fmt.Printf("%s [%s] -> %s (score: %d)\n",
				c.PGNHistory(), best.Name, best.Note, best.Score)
			c.walk(id, depth+1, best.Move)
		}
	}
}

func main() {
	ob, err := OpenOB("./openings.db")
	quitOn(err)
	defer ob.Close()
	coach, err := NewCoach("./stockfish", "./openings.db")
	quitOn(err)
	coach.walk(0, 0, nil)
}
