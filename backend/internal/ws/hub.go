package ws

var hub *Hub = &Hub{
	Clients:    make(map[string]*Client),
	Broadcast:  make(chan Message),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
}

type Hub struct {
	Clients    map[string]*Client
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client
}

func GetHub() *Hub {
	return hub
}

func GetClient(sessionID string) (*Client, bool) {
	client, ok := hub.Clients[sessionID]
	return client, ok
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.SessionID] = client
		case client := <-h.Unregister:
			if _, ok := h.Clients[client.SessionID]; ok {
				delete(h.Clients, client.SessionID)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			for sessionID, client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, sessionID)
				}
			}
		}
	}
}
