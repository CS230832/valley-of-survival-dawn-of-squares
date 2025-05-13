package game

import (
	"log"
	"time"
	"valley-of-survival-dawn-of-squares/internal/ws"
)

var tickRate = 60
var tickInterval = time.Second / time.Duration(tickRate)

func StartWorld() {
	ticker := time.NewTicker(tickInterval)
	defer ticker.Stop()

	for range ticker.C {
		updateWorld()
	}
}

func updateWorld() {
	hub := ws.GetHub()

	select {
	case message := <-hub.Broadcast:
		log.Println(message)
	default:
	}
}
