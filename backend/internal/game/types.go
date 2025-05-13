package game

type (
	User struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Password string `json:"-"`
	}

	Clan struct {
		ID       uint   `json:"id"`
		Name     string `json:"name"`
		Password string `json:"-"`
		OwnerID  uint   `json:"owner_id"`
	}

	Player struct {
		ID          uint    `json:"id"`
		UserID      uint    `json:"user_id"`
		HP          uint    `json:"hp"`
		Position    [2]uint `json:"position"`
		Color       string  `json:"color"`
		Texturepath string  `json:"texture_path,omitempty"`
		EXP         uint    `json:"exp"`
		ClanID      *uint   `json:"clan_id,omitempty"`
	}

	WeaponClass struct {
		ID             uint    `json:"id"`
		BaseDamage     uint    `json:"base_damage"`
		BaseRange      float32 `json:"base_range"`
		BaseRateOfFire uint    `json:"base_rate_of_fire"`
	}

	Weapon struct {
		ID            uint `json:"id"`
		WeaponClassID uint `json:"weapon_class_id"`
		Level         uint `json:"level"`
		PlayerID      uint `json:"player_id"`
	}

	Enemy struct {
		ID         uint    `json:"id"`
		Position   [2]uint `json:"position"`
		Damage     uint    `json:"-"`
		Range      float32 `json:"-"`
		RateOfFire uint    `json:"-"`
		Size       uint    `json:"size"`
	}
)
