package game

import "sync"

const WorldSize = 2048
const HalfWorldSize = 1024
const PlayerSize = 32
const HalfPlayerSize = 16

type GameState struct {
	SpawnedPlayersMu sync.RWMutex     `json:"-"`
	SpawnedPlayers   map[uint]*Player `json:"players"`

	EnemiesMu sync.RWMutex    `json:"-"`
	Enemies   map[uint]*Enemy `json:"enemies"`
}

// var gameState GameState
