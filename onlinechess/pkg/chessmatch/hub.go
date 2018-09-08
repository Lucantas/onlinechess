package chessmatch

type movement struct {
	data  []byte
	match string
}

// Hub represents basic server informations of the websocket
// server, such as clients connected, games, movements.
type Hub struct {
	// Clients registered on the hub
	Clients map[*player]bool
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
		Clients:    make(map[*player]bool),
		register:   make(chan *player),
		unregister: make(chan *player),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case player := <-h.register:
			h.Clients[player] = true
		case player := <-h.unregister:
			if _, ok := h.Clients[player]; ok {
				delete(h.Clients, player)
				close(player.Client.Send)
			}
		case movement := <-h.Movement:
			for player := range h.Clients {
				select {
				case player.Client.Send <- movement.data:
				default:
					close(player.Client.Send)
					delete(h.Clients, player)
				}
			}
		}
	}
}
