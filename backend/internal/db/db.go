package db

import (
	"context"
	"log"
	"valley-of-survival-dawn-of-squares/internal/game"

	"github.com/jackc/pgx/v5"
)

var conn *pgx.Conn

func InitDB() {
	var err error
	conn, err = pgx.Connect(context.Background(), "postgres://appuser:apppassword@localhost:5432/appdb?sslmode=disable")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	log.Println("Hey, Successfully connected to PostgreSQL!")
}

func CloseDB() {
	if err := conn.Close(context.Background()); err != nil {
		log.Fatalf("Failed to close db connection: %v", err)
	}
	log.Println("Successfully closed connection to PostgreSQL!")
}

func GetDB() *pgx.Conn {
	return conn
}

func CreateUser(c context.Context, name string, password string) error {
	tx, err := conn.BeginTx(c, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(c)

	_, err = tx.Exec(c, `
		INSERT INTO users (name, password) VALUES ($1, $2)
	`, name, password)
	if err != nil {
		return err
	}

	return tx.Commit(c)
}

func GetUserByName(c context.Context, username string) (*game.User, error) {
	tx, err := conn.BeginTx(c, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(c)

	row := tx.QueryRow(c, `
		SELECT id, name, password FROM users WHERE name = $1
	`, username)

	var u game.User
	if err := row.Scan(&u.ID, &u.Username, &u.Password); err != nil {
		return nil, err
	}

	if err := tx.Commit(c); err != nil {
		return nil, err
	}

	return &u, nil
}

func GetUserByID(c context.Context, id uint) (*game.User, error) {
	tx, err := conn.BeginTx(c, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(c)

	row := tx.QueryRow(c, `
		SELECT id, name, password FROM users WHERE id = $1
	`, id)

	var u game.User
	if err := row.Scan(&u.ID, &u.Username, &u.Password); err != nil {
		return nil, err
	}

	if err := tx.Commit(c); err != nil {
		return nil, err
	}

	return &u, nil
}

func CreatePlayer(
	c context.Context,
	userID uint,
	hp uint,
	position [2]uint,
	color string,
	texturepath string,
	exp uint,
) error {
	tx, err := conn.Begin(c)
	if err != nil {
		return err
	}
	defer tx.Rollback(c)

	_, err = tx.Exec(c, `
		INSERT INTO players (user_id, hp, position_x, position_y, color, texture_path, experience_points)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, userID, hp, position[0], position[1], color, texturepath, exp)
	if err != nil {
		return err
	}

	return tx.Commit(c)
}

func GetPlayerByUsername(c context.Context, username string) (*game.Player, error) {
	tx, err := conn.BeginTx(c, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(c)

	row := tx.QueryRow(c, `
		SELECT p.id, p.user_id, p.hp, p.position_x, p.position_y, p.color, p.texture_path, p.experience_points, p.clan_id
		FROM players p
		JOIN users u ON u.id = p.user_id
		WHERE u.name = $1
	`, username)

	var p game.Player
	var posX, posY uint
	if err := row.Scan(&p.ID, &p.UserID, &p.HP, &posX, &posY, &p.Color, &p.Texturepath, &p.EXP, &p.ClanID); err != nil {
		return nil, err
	}
	p.Position = [2]uint{posX, posY}

	if err := tx.Commit(c); err != nil {
		return nil, err
	}

	return &p, nil
}

func GetPlayerByID(c context.Context, id uint) (*game.Player, error) {
	tx, err := conn.BeginTx(c, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(c)

	row := tx.QueryRow(c, `
		SELECT p.id, p.user_id, p.hp, p.position_x, p.position_y, p.color, p.texture_path, p.experience_points, p.clan_id
		FROM players p
		WHERE p.id = $1
	`, id)

	var p game.Player
	if err := row.Scan(&p.ID, &p.UserID, &p.HP, &p.Position[0], &p.Position[1], &p.Color, &p.Texturepath, &p.EXP, &p.ClanID); err != nil {
		return nil, err
	}

	if err := tx.Commit(c); err != nil {
		return nil, err
	}

	return &p, nil
}

func GetWeaponClass(c context.Context, id uint) (*game.WeaponClass, error) {
	tx, err := conn.BeginTx(c, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(c)

	row := tx.QueryRow(c, `
		SELECT id, base_damage, base_range, base_rate_of_fire
		FROM weapon_classes
		WHERE id = $1
	`, id)

	var wc game.WeaponClass
	if err := row.Scan(&wc.ID, &wc.BaseDamage, &wc.BaseRange, &wc.BaseRateOfFire); err != nil {
		return nil, err
	}

	if err := tx.Commit(c); err != nil {
		return nil, err
	}

	return &wc, nil
}

func CreateWeapon(
	c context.Context,
	weaponClassID uint,
	level uint,
	playerID uint,
) (*game.Weapon, error) {
	tx, err := conn.BeginTx(c, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(c)

	var id uint
	err = tx.QueryRow(c, `
		INSERT INTO weapons (weapon_class_id, level, player_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`, weaponClassID, level, playerID).Scan(&id)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(c); err != nil {
		return nil, err
	}

	return &game.Weapon{
		ID:            id,
		WeaponClassID: weaponClassID,
		Level:         level,
		PlayerID:      playerID,
	}, nil
}

func GetWeapon(c context.Context, id uint) (*game.Weapon, error) {
	tx, err := conn.BeginTx(c, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(c)

	row := tx.QueryRow(c, `
		SELECT id, weapon_class_id, level, player_id
		FROM weapons
		WHERE id = $1
	`, id)

	var w game.Weapon
	if err := row.Scan(&w.ID, &w.WeaponClassID, &w.Level, &w.PlayerID); err != nil {
		return nil, err
	}

	if err := tx.Commit(c); err != nil {
		return nil, err
	}

	return &w, nil
}

func GetPlayerWeapons(c context.Context, id uint) ([]*game.Weapon, error) {
	tx, err := conn.BeginTx(c, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(c)

	rows, err := tx.Query(c, `
		SELECT id, weapon_class_id, level, player_id
		FROM weapons
		WHERE player_id = $1;
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var weapons []*game.Weapon
	for rows.Next() {
		var w game.Weapon
		if err := rows.Scan(&w.ID, &w.WeaponClassID, &w.Level, &w.PlayerID); err != nil {
			return nil, err
		}
		weapons = append(weapons, &w)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := tx.Commit(c); err != nil {
		return nil, err
	}

	return weapons, nil
}

func GetClanPlayers(c context.Context, id uint) ([]*game.Player, error) {
	tx, err := conn.BeginTx(c, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(c)

	rows, err := tx.Query(c, `
	SELECT p.id, p.user_id, p.hp, p.position_x, p.position_y, p.color, p.texture_path, p.experience_points, p.clan_id
	FROM players p
	WHERE p.clan_id = $1;
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []*game.Player
	for rows.Next() {
		var player game.Player
		if err := rows.Scan(&player.ID, &player.UserID, &player.HP, &player.Position[0], &player.Position[1], &player.Color, &player.Texturepath, &player.EXP, &player.ClanID); err != nil {
			return nil, err
		}
		players = append(players, &player)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := tx.Commit(c); err != nil {
		return nil, err
	}

	return players, nil
}

func CreateClan(c context.Context, name string, password string, ownerID uint) error {
	tx, err := conn.Begin(c)
	if err != nil {
		return err
	}
	defer tx.Rollback(c)

	if _, err = tx.Exec(c, `
		INSERT INTO clans (name, password, owner_id)
		VALUES ($1, $2, $3);
	`, name, password, ownerID); err != nil {
		return err
	}

	return tx.Commit(c)
}

func GetClanByUsername(c context.Context, username string) (*game.Clan, error) {
	var clan game.Clan
	err := conn.QueryRow(c, `
		SELECT clans.id, clans.name, clans.password, clans.owner_id
		FROM clans
		JOIN users u ON u.id = clans.owner_id
		WHERE u.name = $1;
	`, username).Scan(&clan.ID, &clan.Name, &clan.Password, &clan.OwnerID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // No clan found
		}
		return nil, err
	}

	return &clan, nil
}

func GetClanByID(c context.Context, id uint) (*game.Clan, error) {
	var clan game.Clan
	err := conn.QueryRow(c, `
		SELECT id, name, password, owner_id
		FROM clans
		WHERE id = $1
	`, id).Scan(&clan.ID, &clan.Name, &clan.Password, &clan.OwnerID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // No clan found
		}
		return nil, err
	}

	return &clan, nil
}

func GetClanByName(c context.Context, name string) (*game.Clan, error) {
	var clan game.Clan
	err := conn.QueryRow(c, `
		SELECT id, name, password, owner_id
		FROM clans
		WHERE name = $1
	`, name).Scan(&clan.ID, &clan.Name, &clan.Password, &clan.OwnerID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // No clan found
		}
		return nil, err
	}

	return &clan, nil
}

func DeleteClan(c context.Context, id uint) error {
	tx, err := conn.Begin(c)
	if err != nil {
		return err
	}
	defer tx.Rollback(c)

	// Delete the clan from the database
	if _, err = tx.Exec(c, `
		DELETE FROM clans WHERE id = $1
	`, id); err != nil {
		return err
	}

	if _, err = tx.Exec(c, `
		UPDATE players
		SET clan_id = NULL
		WHERE clan_id = $1
	`, id); err != nil {
		return err
	}

	return tx.Commit(c)
}

func JoinClan(c context.Context, userID uint, clanID uint) error {
	tx, err := conn.Begin(c)
	if err != nil {
		return err
	}
	defer tx.Rollback(c)

	_, err = tx.Exec(c, `
	UPDATE players
	SET clan_id = $1
	WHERE user_id = $2;
	`, clanID, userID)

	if err != nil {
		return err
	}

	return tx.Commit(c)
}

func LeaveClan(c context.Context, userID uint) error {
	tx, err := conn.Begin(c)
	if err != nil {
		return err
	}
	defer tx.Rollback(c)

	_, err = tx.Exec(c, `
	UPDATE players
	SET clan_id = NULL
	WHERE user_id = $1;
	`, userID)

	if err != nil {
		return err
	}

	return tx.Commit(c)
}

func SpawnPlayer(c context.Context, id uint) error {
	tx, err := conn.Begin(c)
	if err != nil {
		return err
	}
	defer tx.Rollback(c)

	_, err = tx.Exec(c, `
		UPDATE players
		SET is_spawned = TRUE
		WHERE id = $1
	`, id)
	if err != nil {
		return err
	}

	return tx.Commit(c)
}

func DespawnPlayer(c context.Context, id uint) error {
	tx, err := conn.Begin(c)
	if err != nil {
		return err
	}
	defer tx.Rollback(c)

	_, err = tx.Exec(c, `
		UPDATE players
		SET is_spawned = FALSE
		WHERE id = $1
	`, id)
	if err != nil {
		return err
	}

	return tx.Commit(c)
}
