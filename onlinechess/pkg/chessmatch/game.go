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

type player struct {
	Client *Client
	match  string
}

func (p *player) Read() {
	// Starts webscoket connection
	conn := p.Client.Conn
	defer func() {
		// close the connection after the end of function
		p.Client.Hub.unregister <- p
		conn.Close()
	}()
	// loop on the data
	for {
		_, move, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("Error: ", err)
			}
			break
		}
		// process and pass the data trought the channel to the player
		m := movement{move, p.match}
		p.Client.Hub.Movement <- m
	}

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

	go player.Write()
	go player.Read()

}

func newPlayer(client *Client, match string) *player {
	return &player{client, match}
}
