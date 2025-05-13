package session

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"sync"
)

type Session struct {
	SessionID string
	Username  string
}

var sessions = make(map[string]string)
var mu sync.RWMutex

func CreateSession(username string) string {
	mu.Lock()
	defer mu.Unlock()

	b := make([]byte, 16)
	rand.Read(b)
	sessionID := hex.EncodeToString(b)
	sessions[sessionID] = username

	log.Println("created a session:", sessionID, "->", username)

	return sessionID
}

func GetUsername(sessionID string) (string, bool) {
	mu.RLock()
	defer mu.RUnlock()

	username, ok := sessions[sessionID]
	return username, ok
}

func RemoveSession(sessionID string) {
	mu.Lock()
	defer mu.Unlock()

	if username, ok := sessions[sessionID]; ok {
		log.Println("removed a session:", sessionID, "->", username)

		delete(sessions, sessionID)
	}
}
