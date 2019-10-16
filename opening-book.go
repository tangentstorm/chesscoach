package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/mattn/go-sqlite3"
	"github.com/notnil/chess"
)

func catch(err error) {
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
        prev integer not null default 0, move text );`)
		if err == nil {
			ob = &OpeningBook { db }
		}
	}
	return ob, err
}

func (ob *OpeningBook) Line(n int) (game *chess.Game, err error) {
	stmt, err := ob.db.Prepare(`
	with recursive line(pre) as
    (values(?) union select prev from move, line where id=line.pre)
  select id, move from move where id in line`)
	catch(err)
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


func main() {
	ob, err := OpenOB("./openings.db")
	catch(err)
	defer ob.Close()
	game, err := ob.Line(2)
	catch(err)
	fmt.Println(game)
	fmt.Println(game.Position().Board().Draw())
}
