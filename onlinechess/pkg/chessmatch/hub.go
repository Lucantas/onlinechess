package chessmatch

import (
	"encoding/json"
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
	//Games map[string]map[*player]bool
	Games []*match
	// movement made by a Client
	Messages chan []byte
	// Register requests from the clients.
	register chan *player
	// Unregister requests from clients.
	unregister chan *player
}

// NewHub initialize and returns a pointer to a Hub
func NewHub() *Hub {
	return &Hub{
		Messages:   make(chan []byte),
		Lobby:      make(map[*player]bool),
		register:   make(chan *player),
		unregister: make(chan *player),
	}
}

func matchInfo() {
	in := `{"matchInfo":"{
			"status":"playing",
			"player1":p1,
			"player2":p2,
			"whites":"player1",
			"blacks":"player2",
		}"
	}`

	rawIn := json.RawMessage(in)
	bytes, err := rawIn.MarshalJSON()
	if err != nil {
		panic(err)
	}

	var p Person
	err = json.Unmarshal(bytes, &p)
	if err != nil {
		panic(err)
	}
}

func matchMaker(h *Hub, player *player) {
	for _, match := range h.Games {
		log.Println("test cookie:", player.Client.Cookie, match.Player1.Client.Cookie, match.Player2.Client.Cookie)
		if player.Client.Cookie == match.Player1.Client.Cookie || player.Client.Cookie == match.Player2.Client.Cookie {
			log.Println("user have match")
			player.Client.Send <- []byte("Match Info:")
		}
	}
	log.Println(player.match.ID)
	for p := range h.Lobby {
		if h.Lobby[p] && p != player {
			rnd := rand.New(rand.NewSource(int64(p.Client.ID)))
			m := newMatch(rnd.Int(), p, player)
			p.match, player.match = m, m
			p.Client.Send <- []byte("Match Found")
			player.Client.Send <- []byte("Match Found")
			log.Println("player cookie", player.Client.Cookie)
			log.Println("player2 cookie", m.Player2.Client.Cookie)
			h.Lobby[p] = false
			h.Lobby[player] = false
			if len(h.Games) > 0 {
				for _, match := range h.Games {
					if match.ID != m.ID {
						h.Games = append(h.Games, m)
					}
				}
			} else {
				h.Games = append(h.Games, m)
			}
			log.Println("length of games on hub after append:", len(h.Games))
		}
	}
}

func (h *Hub) Run() {
	for {
		select {
		case player := <-h.register:
			log.Println("Player registering on lobby")
			h.Lobby[player] = true
			matchMaker(h, player)
		case player := <-h.unregister:
			if _, ok := h.Lobby[player]; ok {
				delete(h.Lobby, player)
				close(player.Client.Send)
			}
		case message := <-h.Messages:
			for player := range h.Lobby {
				select {
				case player.Client.Send <- message:
				default:
					close(player.Client.Send)
					delete(h.Lobby, player)
				}
			}

		}
	}
}
