package utils

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"valley-of-survival-dawn-of-squares/internal/session"
)

var clientSessionMu sync.RWMutex
var clientSessions = make(map[string]string)
var sessionClients = make(map[string]string)

func GetSession(r *http.Request) (*session.Session, bool) {
	session, ok := r.Context().Value(session.Session{}).(*session.Session)
	return session, ok
}

func GetClientIdentifier(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}

	ip = strings.Split(ip, ",")[0]
	ip = strings.Split(ip, ":")[0]

	return fmt.Sprintf("IP: %s | User-Agent: %s", ip, r.UserAgent())
}

func SetClientSession(sessionID string, r *http.Request) {
	clientSessionMu.Lock()
	defer clientSessionMu.Unlock()
	clientSessions[sessionID] = GetClientIdentifier(r)
	sessionClients[GetClientIdentifier(r)] = sessionID
}

func RemoveClientSession(sessionID string) {
	clientSessionMu.Lock()
	defer clientSessionMu.Unlock()
	delete(clientSessions, sessionID)
	if client, ok := sessionClients[sessionID]; ok {
		delete(sessionClients, client)
	}
}

func GetClientSession(sessionID string) (string, bool) {
	clientSessionMu.RLock()
	defer clientSessionMu.RUnlock()
	client, ok := clientSessions[sessionID]
	return client, ok
}
