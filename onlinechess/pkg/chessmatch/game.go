package chessmatch

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  528,
	WriteBufferSize: 528,
}

// Client represents a client who asked for a connection
type Client struct {
	Hub *Hub

	//the match name
	Match string

	// websocket connection
	Conn *websocket.Conn

	// channel of outbund movements
	Send chan []byte
}

func (p *player) Read() {
	// Starts webscoket connection

	// close of the connection after the end of function

	// loop on the data

	// pass the data trought the channel to the player
}

func (p *player) Write() {
	// Starts a websocket connection

	// process the data to be sent over the socket

	// send the data to the player's room
}

func socketHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(conn)
	client := &Client{Hub: hub, Conn: conn, Send: make(chan []byte, 256)}
	player := newPlayer(client, "match")
	client.Hub.register <- player
	//match :=  r.URL.Query()["match"][0]

	go player.Read()
	go player.Write()

}

func newPlayer(client *Client, match string) *player {
	return &player{client, match}
}
