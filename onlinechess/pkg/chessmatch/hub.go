package chessmatch

import (
	"log"
	"math/rand"
)

type movement struct {
	data  []byte
	match *match
}

// Hub represents basic server informations of the websocket
// server, such as clients connected, games, movements.
type Hub struct {
	// Lobby holds the clients registered on the hub
	Lobby map[*player]bool
	// Games registered within the clients
	Games map[string]map[*player]bool
	// movement made by a Client
	Movement chan movement
	// Register requests from the clients.
	register chan *player
	// Unregister requests from clients.
	unregister chan *player
}

// NewHub initialize and returns a pointer to a Hub
func NewHub() *Hub {
	return &Hub{
		Movement:   make(chan movement),
		Lobby:      make(map[*player]bool),
		register:   make(chan *player),
		unregister: make(chan *player),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case player := <-h.register:
			h.Lobby[player] = true
			for p := range h.Lobby {
				if h.Lobby[p] && p != player {
					log.Println("Found player on lobby")
					rnd := rand.New(rand.NewSource(int64(p.Client.ID)))
					m := newMatch(rnd.Int(), p, player)
					p.match, player.match = m, m
					p.Client.Send <- []byte("Match Found")
					player.Client.Send <- []byte("Match Found")
					h.Lobby[p] = false
					h.Lobby[player] = false
				}
			}
		case player := <-h.unregister:
			if _, ok := h.Lobby[player]; ok {
				delete(h.Lobby, player)
				close(player.Client.Send)
			}
		case movement := <-h.Movement:
			for player := range h.Lobby {
				select {
				case player.Client.Send <- movement.data:
				default:
					close(player.Client.Send)
					delete(h.Lobby, player)
				}
			}
		}
	}
}
