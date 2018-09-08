package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"websockets/onlinechess/pkg/chessmatch"
)

func main() {
	hub := chessmatch.NewHub()
	go hub.Run()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./public/static/"))))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		chessmatch.SocketHandler(hub, w, r)
	})
	addr := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(addr, nil))

}
