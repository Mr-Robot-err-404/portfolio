package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {
	server := Server{upgrader: &websocket.Upgrader{}}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/", http.FileServer(http.Dir("web")))

	http.HandleFunc("/shell", server.handle)

	log.Println("portfolio available at http://localhost:4242")
	log.Fatal(http.ListenAndServe(":4242", nil))
}
