package chessmatch

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  528,
		WriteBufferSize: 528,
	}
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client represents a client who asked for a connection
type Client struct {
	ID  uint
	Hub *Hub

	//the match name
	match string

	// websocket connection
	Conn *websocket.Conn

	// channel of outbund movements
	Send chan []byte
}

type player struct {
	Client *Client
	match  *match
}

// match is the actual game, aggregating two players, a name and a Id
type match struct {
	ID      uint
	Name    string
	Player1 *player
	Player2 *player
}

func (p *player) FindMatch() {
	// Starts webscoket connection
	conn := p.Client.Conn
	defer func() {
		// close the connection after the end of function
		p.Client.Hub.unregister <- p
		conn.Close()
	}()
	for {
		_, message, err := conn.ReadMessage()
		log.Println("message:", message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("Error: ", err)
			}
			break
		}
		// process and pass the data trought the channel to the player
		m := movement{message, p.match}
		p.Client.Hub.Movement <- m
	}

}

func (p *player) Read() {
	// Starts webscoket connection
	conn := p.Client.Conn
	defer func() {
		// close the connection after the end of function
		p.Client.Hub.unregister <- p
		conn.Close()
	}()
	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
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
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		// close the connection after the end of function
		ticker.Stop()
		conn.Close()
	}()

	for {
		log.Println("on write")
		select {
		case movement, ok := <-p.Client.Send:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
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
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
	// process the data to be sent over the socket

	// send the data to the player's room
}

func newMatch(id int, p1 *player, p2 *player) *match {
	m := &match{
		ID:      uint(id),
		Player1: p1,
		Player2: p2,
	}
	return m
}

// GameSocket handles the websocket connection between two clients
// to provide a private channel for a match
func GameSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	player := newPlayer(hub, conn, "match")
	player.Client.Hub.register <- player
	//match :=  r.URL.Query()["match"][0]

	go player.Write()
	go player.Read()

}

// LobbySocket handles the connection between the client and the
// lobby on the server, enabling the user to find an opponent
func LobbySocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	player := newPlayer(hub, conn, "match")
	player.Client.Hub.register <- player

	go player.Write()
	go player.FindMatch()
}

func newPlayer(hub *Hub, conn *websocket.Conn, matchName string) *player {
	client := &Client{Hub: hub, Conn: conn, Send: make(chan []byte, 256)}
	return &player{client, &match{Name: matchName}}
}
