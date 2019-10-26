package main

// Mini HTTP + Websocket server for chesscoach.
// (The websockets don't do anything yet.)

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {
	r := http.NewServeMux()
	r.Handle("/sprites/", http.FileServer(http.Dir("."))) // url not stripped, so it serves ./sprites/
	r.Handle("/", http.FileServer(http.Dir("./htmlui")))
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		up := websocket.Upgrader{}
		ws, err := up.Upgrade(w, r, nil)
		log.Fatal(err)
		for {
			msgType, msg, err := ws.ReadMessage()
			log.Fatal(err)
			fmt.Println(string(msg))
			// echo the message back to the socket:
			err = ws.WriteMessage(msgType, msg)
			log.Fatal(err)
		}
	})
	log.Fatal(http.ListenAndServe(":8080", r))
}
