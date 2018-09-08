package chessmatch

import (
	"log"
	"net/http"
	"time"

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
	// Starts webscoket connection
	conn := p.Client.Conn
	ticker := time.NewTicker(15)
	defer func() {
		// close the connection after the end of function
		ticker.Stop()
		conn.Close()
	}()

	for {
		select {
		case movement, ok := <-p.Client.Send:
			if !ok {
				// The hub closed the channel.
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			log.Println(movement)
			w.Write(movement)

			// Add queued chat messages to the current websocket message.
			n := len(p.Client.Send)
			for i := 0; i < n; i++ {
				w.Write(<-p.Client.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
	// process the data to be sent over the socket

	// send the data to the player's room
}

// SocketHandler handles the websocket connection between the client
// and the server, trough the net/http package
func SocketHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
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
