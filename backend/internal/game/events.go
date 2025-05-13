package game

// Client

var (
	ClientKeyPressed    = "ClientKeyPressed"    // data -> Key
	ClientKeyDown       = "ClientKeyDown"       // data -> Key
	ClientPlayerSpawn   = "ClientPlayerSpawn"   // empty request
	ClientPlayerDespawn = "ClientPlayerDespawn" // empty request
	ClientSelect        = "ClientSelect"        // data -> selection number out of one, two, or three
)

var (
	KeyW = 'w'
	KeyA = 'a'
	KeyS = 's'
	KeyD = 'd'

	One   = '1'
	Two   = '2'
	Three = '3'
)

// Server
type (
	MovedEntity struct {
		ID       uint    `json:"id"`
		Position [2]uint `json:"position"`
	}

	PlayerUpgrade struct {
		ID             uint    `json:"id"`
		WeaponClassIDs [3]uint `json:"weapon_class_ids"`
	}
)

var (
	ServerPlayersSpawn   = "ServerPlayersSpawn"   // data -> []Player
	ServerPlayersDespawn = "ServerPlayersDespawn" // data -> []PlayerID
	ServerEnemiesSpawn   = "ServerEnemySpawn"     // data -> []Enemy
	ServerEnemiesDespawn = "ServerEnemyDespawn"   // data -> []EnemyID
	ServerMovePlayers    = "ServerMovePlayers"    // data -> []MovedEntity
	ServerMoveEnemies    = "ServerMoveEnemies"    // data -> []MovedEntity
	ServerDamagePlayers  = "ServerDamagePlayers"  // data -> []PlayerID
	ServerDamageEnemies  = "ServerDamageEnemies"  // data -> []EnemyID
	ServerUpgradePlayer  = "ServerUpgradePlayer"  // data -> []PlayerUpgrade
)
