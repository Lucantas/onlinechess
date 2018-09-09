package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"websockets/onlinechess/pkg/chessmatch"
)

var (
	address = fmt.Sprintf("http://%s:%s/", os.Getenv("HOST"), os.Getenv("PORT"))
	tpl     = &template.Template{}
)

func init() {
	tpl = template.Must(template.ParseGlob("./public/static/*.gohtml"))
}

func main() {
	hub := chessmatch.NewHub()
	go hub.Run()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./public/static/"))))
	http.HandleFunc("/", home)
	http.HandleFunc("/guest/match", match)
	http.HandleFunc("/guest/register", registerGuest)
	http.HandleFunc("/guest/play", handleGuest)
	http.HandleFunc("/matchws", func(w http.ResponseWriter, r *http.Request) {
		chessmatch.GameSocket(hub, w, r)
	})
	http.HandleFunc("/lobbyws", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			log.Println("no cookie")
			http.Redirect(w, r, "/guest/play", http.StatusSeeOther)
		}
		chessmatch.LobbySocket(hub, cookie, w, r)
	})
	addr := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(addr, nil))

}

func home(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := tpl.ExecuteTemplate(w, "home.gohtml", address); err != nil {
		log.Println("Error: ", err)
	}
	//http.ServeFile(w, r, "./public/static/home.html")
}

func match(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/guest/play", http.StatusSeeOther)
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := tpl.ExecuteTemplate(w, "chess.gohtml", address); err != nil {
		log.Println("Error: ", err)
	}
}

func registerGuest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user := r.FormValue("username")
	if user == " " || user == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil {
		newSession(w, user)
	}
	log.Println(cookie)
	http.Redirect(w, r, "/guest/play", http.StatusSeeOther)
}

func handleGuest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Test if user got cookies
	if _, err := r.Cookie("session"); err != nil {
		// if user have no cookie defined, open the guest template to input the username
		if err := tpl.ExecuteTemplate(w, "guest.gohtml", address); err != nil {
			log.Println("Error: ", err)
		}
	} else {
		// if user have a cookie, open the lobby template while a match is being searched
		if err := tpl.ExecuteTemplate(w, "lobby.gohtml", address); err != nil {
			log.Println("Error: ", err)
		}
	}
}
