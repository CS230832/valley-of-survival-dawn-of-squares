package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	writeWait  = 1 * time.Second
	pongWait   = 5 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type Client struct {
	Hub       *Hub
	Conn      *websocket.Conn
	SessionID string
	Send      chan Message
	Read      chan Message
}

type Message struct {
	SessionID string `json:"session"`
	Type      string `json:"type"`
	Data      any    `json:"data"`
}

func newClient(conn *websocket.Conn, sessionID string) *Client {
	return &Client{
		Hub:       hub,
		Conn:      conn,
		SessionID: sessionID,
		Send:      make(chan Message),
		Read:      hub.Broadcast,
	}
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session")
	if len(sessionID) == 0 {
		http.Error(w, "no session query paramter given", http.StatusBadRequest)
		return
	}

	wsconn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		http.Error(w, "Upgrade failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	client := newClient(wsconn, sessionID)
	client.Hub.Register <- client

	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		var msg Message
		if err := c.Conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("Unexpected close error:", err)
			}
			break
		}
		c.Hub.Broadcast <- msg
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Println("Error getting writer:", err)
				return
			}

			jsonData, err := json.Marshal(message)
			if err != nil {
				log.Println("Error marshaling JSON:", err)
				return
			}

			w.Write(jsonData)

			for range len(c.Send) {
				if _, err := w.Write([]byte("\n")); err != nil {
					log.Println("Error writing newline:", err)
					return
				}

				queued := <-c.Send
				jsonQueued, err := json.Marshal(queued)
				if err != nil {
					log.Println("Error marshaling queued message:", err)
					return
				}

				if _, err := w.Write(jsonQueued); err != nil {
					log.Println("Error writing queued message:", err)
					return
				}
			}

			if err := w.Close(); err != nil {
				log.Println("Error closing writer:", err)
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("Error writing ping:", err)
				return
			}
		}
	}
}
