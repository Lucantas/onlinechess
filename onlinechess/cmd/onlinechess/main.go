package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"websockets/onlinechess/pkg/chessmatch"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "./public/static/home.html")
}

func serveMatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "./public/static/chess.html")
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(dir)
	hub := chessmatch.NewHub()
	go hub.Run()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./public/static/"))))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		chessmatch.SocketHandler(hub, w, r)
	})
	addr := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(addr, nil))

}
