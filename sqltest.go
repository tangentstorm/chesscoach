package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func catch(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	db, err := sql.Open("sqlite3", "./openings.db")
	catch(err)
	defer db.Close()

	_, err = db.Exec(`
	create table if not exists position (
    id integer not null primary key autoincrement,
    board text );

  -- starting position
  insert or ignore into position (id, board) values
    (0, 'rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR');

  create table if not exists move (
    id integer primary key autoincrement,
    prev integer not null default 0, move text );

  insert or ignore into move (id, prev, move) values (1, 0, 'e4');
  insert or ignore into move (id, prev, move) values (2, 1, 'e5');

  `)
	catch(err)


	stmt, err := db.Prepare(`
	with recursive line(pre) as
    (values(?) union select prev from move, line where id=line.pre)
  select id, move from move where id in line`)
	catch(err)
	defer stmt.Close()

	rows, err := stmt.Query("2") // query the line for move with id=2
	catch(err)
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		catch(err)
		fmt.Println(id, name)
	}
	catch(rows.Err())
}
