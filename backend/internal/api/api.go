package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"valley-of-survival-dawn-of-squares/internal/db"
	"valley-of-survival-dawn-of-squares/internal/game"
	"valley-of-survival-dawn-of-squares/internal/session"
	"valley-of-survival-dawn-of-squares/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ClanPayload struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

var VosDosSessionToken = "vosdos-session-token"

func HandlerWithAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.Header.Get(VosDosSessionToken)
		if len(sessionID) == 0 {
			http.Error(w, "user not logged in", http.StatusBadRequest)
			return
		}

		username, exists := session.GetUsername(sessionID)
		if !exists {
			http.Error(w, "invalid session token", http.StatusBadRequest)
			return
		}

		clientId, exists := utils.GetClientSession(sessionID)

		if !exists || strings.Compare(strings.TrimSpace(clientId), strings.TrimSpace(utils.GetClientIdentifier(r))) != 0 {
			http.Error(w, "invalid session token", http.StatusBadRequest)
			return
		}

		handlerFunc(w, r.WithContext(
			context.WithValue(
				r.Context(),
				session.Session{},
				&session.Session{SessionID: sessionID, Username: username},
			),
		))
	}
}

func HandleVerifySession(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func HandleSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := db.GetUserByName(r.Context(), creds.Username); err == nil {
		http.Error(w, "user already exists", http.StatusConflict)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := db.CreateUser(r.Context(), creds.Username, string(hashedPassword)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user, err := db.GetUserByName(r.Context(), creds.Username); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if err = db.CreatePlayer(r.Context(), user.ID, 100, [2]uint{game.HalfWorldSize, game.HalfWorldSize}, "white", "", 0); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user, err := db.GetUserByName(r.Context(), creds.Username); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)) != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if sessionID := r.Header.Get(VosDosSessionToken); len(sessionID) != 0 {
		session.RemoveSession(sessionID)
		utils.RemoveClientSession(sessionID)
	}

	sessionID := session.CreateSession(creds.Username)
	utils.SetClientSession(sessionID, r)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(sessionID))
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get(VosDosSessionToken)

	session.RemoveSession(sessionID)
	utils.RemoveClientSession(sessionID)

	w.WriteHeader(http.StatusOK)
}

func HandleGetUserInfo(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	var user *game.User
	var err error

	userIDStr := r.URL.Query().Get("user_id")

	if len(userIDStr) != 0 {
		userID, err := strconv.Atoi(userIDStr)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user, err = db.GetUserByID(r.Context(), uint(userID))

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if username := queryParams.Get("username"); len(username) != 0 {
		user, err = db.GetUserByName(r.Context(), username)
	} else {
		http.Error(w, "no username nor user_id query paramter given", http.StatusBadRequest)
		return

	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userJson, err := json.Marshal(*user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userJson)
}

func HandleGetPlayerInfo(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	var player *game.Player
	var err error

	playerIDStr := r.URL.Query().Get("player_id")

	if len(playerIDStr) != 0 {
		playerID, err := strconv.Atoi(playerIDStr)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		player, err = db.GetPlayerByID(r.Context(), uint(playerID))

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if username := queryParams.Get("username"); len(username) != 0 {
		player, err = db.GetPlayerByUsername(r.Context(), username)
	} else {
		http.Error(w, "no username nor player_id query paramter given", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	playerJson, err := json.Marshal(*player)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(playerJson)
}

func HandleGetClanInfo(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	var clan *game.Clan
	var err error

	if clandIDStr := queryParams.Get("clan_id"); len(clandIDStr) != 0 {
		clanID, err := strconv.Atoi(clandIDStr)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		clan, err = db.GetClanByID(r.Context(), uint(clanID))

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if clanName := r.URL.Query().Get("clan_name"); len(clanName) != 0 {
		clan, err = db.GetClanByName(r.Context(), clanName)
	} else {
		http.Error(w, "no clan_name query paramter given", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	clanJson, err := json.Marshal(*clan)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(clanJson)
}

func HandleGetWeaponClassInfo(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	var weaponClass *game.WeaponClass
	var err error

	weaponClassIDStr := queryParams.Get("weapon_class_id")

	if len(weaponClassIDStr) != 0 {
		weaponClassID, err := strconv.Atoi(weaponClassIDStr)

		if err != nil {
			http.Error(w, "invalid weapon_class_id", http.StatusBadRequest)
			return
		}

		weaponClass, err = db.GetWeaponClass(r.Context(), uint(weaponClassID))

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "no weapon_class_id query parameter given", http.StatusBadRequest)
		return
	}

	weaponClassJson, err := json.Marshal(*weaponClass)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(weaponClassJson)
}

func HandleGetWeaponInfo(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	var weapon *game.Weapon
	var err error

	weaponIDStr := queryParams.Get("weapon_id")

	if len(weaponIDStr) != 0 {
		weaponID, err := strconv.Atoi(weaponIDStr)

		if err != nil {
			http.Error(w, "invalid weapon_id", http.StatusBadRequest)
			return
		}

		weapon, err = db.GetWeapon(r.Context(), uint(weaponID))

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "no weapon_id query parameter given", http.StatusBadRequest)
		return
	}

	weaponJson, err := json.Marshal(*weapon)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(weaponJson)
}

func HandleGetCurrentUserInfo(w http.ResponseWriter, r *http.Request) {
	session, ok := utils.GetSession(r)
	if !ok {
		http.Error(w, "user not logged in", http.StatusBadRequest)
		return
	}

	user, err := db.GetUserByName(r.Context(), session.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userJson, err := json.Marshal(*user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userJson)
}

func HandleGetCurrentPlayerInfo(w http.ResponseWriter, r *http.Request) {
	session, ok := utils.GetSession(r)
	if !ok {
		http.Error(w, "user not logged in", http.StatusBadRequest)
		return
	}

	player, err := db.GetPlayerByUsername(r.Context(), session.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	playerJson, err := json.Marshal(*player)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(playerJson)
}

func HandleGetCurrentClanInfo(w http.ResponseWriter, r *http.Request) {
	session, ok := utils.GetSession(r)
	if !ok {
		http.Error(w, "user not logged in", http.StatusBadRequest)
		return
	}

	clan, err := db.GetClanByUsername(r.Context(), session.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusExpectationFailed)
		return
	}

	if clan == nil {
		log.Printf("Warning: no clan found for user %s", session.Username)
		http.Error(w, "clan not found", http.StatusNotFound)
		return
	}

	clanJson, err := json.Marshal(*clan)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(clanJson)
}

func HandleGetCurrentWeaponsInfo(w http.ResponseWriter, r *http.Request) {
	session, ok := utils.GetSession(r)
	if !ok {
		http.Error(w, "user not logged in", http.StatusBadRequest)
		return
	}

	player, err := db.GetPlayerByUsername(r.Context(), session.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	weapons, err := db.GetPlayerWeapons(r.Context(), player.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	weaponsJson, err := json.Marshal(weapons)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(weaponsJson)
}

func HandleCreateClan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionToken, ok := utils.GetSession(r)
	if !ok {
		http.Error(w, "user not logged in", http.StatusBadRequest)
		return
	}

	user, err := db.GetUserByName(r.Context(), sessionToken.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	player, err := db.GetPlayerByUsername(r.Context(), user.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if player.ClanID != nil {
		http.Error(w, "user is already in a clan", http.StatusBadRequest)
		return
	}

	var clanPayload ClanPayload
	if err := json.NewDecoder(r.Body).Decode(&clanPayload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.CreateClan(r.Context(), clanPayload.Name, clanPayload.Password, user.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	clan, err := db.GetClanByName(r.Context(), clanPayload.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := db.JoinClan(r.Context(), user.ID, clan.ID); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func HandleDeleteClan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionToken, ok := utils.GetSession(r)
	if !ok {
		http.Error(w, "user not logged in", http.StatusBadRequest)
		return
	}

	user, err := db.GetUserByName(r.Context(), sessionToken.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	clan, err := db.GetClanByUsername(r.Context(), user.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if clan.OwnerID != user.ID {
		http.Error(w, "user is not the owner of the clan", http.StatusBadRequest)
		return
	}

	if err := db.DeleteClan(r.Context(), clan.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func HandleJoinClan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var clanPayload ClanPayload
	if err := json.NewDecoder(r.Body).Decode(&clanPayload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sessionToken, ok := utils.GetSession(r)
	if !ok {
		http.Error(w, "user not logged in", http.StatusBadRequest)
		return
	}

	user, err := db.GetUserByName(r.Context(), sessionToken.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	player, err := db.GetPlayerByUsername(r.Context(), user.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if player.ClanID != nil {
		http.Error(w, "user is already in a clan", http.StatusBadRequest)
		return
	}

	clan, err := db.GetClanByName(r.Context(), clanPayload.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if clan.Password != clanPayload.Password {
		http.Error(w, "invalid password", http.StatusBadRequest)
		return
	}

	if err := db.JoinClan(r.Context(), user.ID, clan.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func HandleLeaveClan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionToken, ok := utils.GetSession(r)
	if !ok {
		http.Error(w, "user not logged in", http.StatusBadRequest)
		return
	}

	user, err := db.GetUserByName(r.Context(), sessionToken.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	player, err := db.GetPlayerByUsername(r.Context(), user.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if player.ClanID == nil {
		http.Error(w, "user is not in a clan", http.StatusBadRequest)
		return
	}

	clan, err := db.GetClanByUsername(r.Context(), user.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if clan.OwnerID == user.ID {
		http.Error(w, "user is the owner of the clan", http.StatusBadRequest)
		return
	}

	if err := db.LeaveClan(r.Context(), user.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
