package main

import (
	"log"
	"net/http"
	"valley-of-survival-dawn-of-squares/internal/api"
	"valley-of-survival-dawn-of-squares/internal/db"
	"valley-of-survival-dawn-of-squares/internal/game"
	"valley-of-survival-dawn-of-squares/internal/ws"
)

func main() {
	db.InitDB()
	defer db.CloseDB()

	http.HandleFunc("/", api.HandleFrontend)

	http.Handle("/api/verify_session", api.HandlerWithAuth(api.HandleVerifySession))

	http.HandleFunc("/api/signup", api.HandleSignup)
	http.HandleFunc("/api/login", api.HandleLogin)
	http.HandleFunc("/api/logout", api.HandleLogout)
	
	http.HandleFunc("/api/info/user", api.HandleGetUserInfo)
	http.HandleFunc("/api/info/player", api.HandleGetPlayerInfo)
	http.HandleFunc("/api/info/clan", api.HandleGetClanInfo)
	http.HandleFunc("/api/info/weapon_class", api.HandleGetWeaponClassInfo)
	http.HandleFunc("/api/info/weapon", api.HandleGetWeaponInfo)

	http.HandleFunc("/api/info/current_user", api.HandlerWithAuth(api.HandleGetCurrentUserInfo))
	http.HandleFunc("/api/info/current_player", api.HandlerWithAuth(api.HandleGetCurrentPlayerInfo))
	http.HandleFunc("/api/info/current_clan", api.HandlerWithAuth(api.HandleGetCurrentClanInfo))
	http.HandleFunc("/api/info/current_weapons", api.HandlerWithAuth(api.HandleGetCurrentWeaponsInfo))

	http.HandleFunc("/api/clan/create", api.HandlerWithAuth(api.HandleCreateClan))
	http.HandleFunc("/api/clan/delete", api.HandlerWithAuth(api.HandleDeleteClan))
	http.HandleFunc("/api/clan/join", api.HandlerWithAuth(api.HandleJoinClan))
	http.HandleFunc("/api/clan/leave", api.HandlerWithAuth(api.HandleLeaveClan))

	http.HandleFunc("/ws", ws.HandleWebSocket)

	go ws.GetHub().Run()
	go game.StartWorld()

	log.Println("Server listening on 0.0.0.0:8080")
	log.Fatalln(http.ListenAndServe("0.0.0.0:8080", nil))
}
