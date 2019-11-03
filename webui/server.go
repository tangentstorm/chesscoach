package main

// Mini HTTP + Websocket server for chesscoach.
// (The websockets don't do anything yet.)

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func quitOn(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	r := http.NewServeMux()
	r.Handle("/sprites/", http.FileServer(http.Dir("."))) // url not stripped, so it serves ./sprites/
	r.Handle("/", http.FileServer(http.Dir("./webui")))
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		up := websocket.Upgrader{}
		ws, err := up.Upgrade(w, r, nil)
		quitOn(err)
		for {
			msgType, msg, err := ws.ReadMessage()
			quitOn(err)
			fmt.Println(string(msg))
			// echo the message back to the socket:
			err = ws.WriteMessage(msgType, msg)
			quitOn(err)
		}
	})
	log.Fatal(http.ListenAndServe(":8080", r))
}
