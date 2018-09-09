package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/satori/go.uuid"
)

func newSession(w http.ResponseWriter, username string) *http.Cookie {
	id, err := uuid.NewV4()
	if err != nil {
		log.Println("Error: ", err)
	}
	cookie := &http.Cookie{
		Name:     "session",
		Value:    id.String(),
		HttpOnly: false,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 72),
	}
	fmt.Printf("Cookie string: %s, Id string: %s, cookie: %v", cookie.String(), id.String(), cookie)
	http.SetCookie(w, cookie)
	return cookie
}
