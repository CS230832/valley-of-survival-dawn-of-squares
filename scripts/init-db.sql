DROP TABLE IF EXISTS weapons;
DROP TABLE IF EXISTS weapon_classes;
DROP TABLE IF EXISTS players;
DROP TABLE IF EXISTS clans;
DROP TABLE IF EXISTS users;
CREATE TABLE users (
    id SERIAL,
    name VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    PRIMARY KEY(id),
    UNIQUE(name)
);
CREATE TABLE clans (
    id SERIAL,
    name VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    owner_id INTEGER NOT NULL,
    PRIMARY KEY(id),
    FOREIGN KEY(owner_id) REFERENCES users(id),
    UNIQUE(name)
);
CREATE TABLE players (
    id SERIAL,
    user_id INTEGER,
    is_spawned BOOLEAN NOT NULL DEFAULT FALSE,
    hp INTEGER NOT NULL DEFAULT 100,
    position_x FLOAT NOT NULL DEFAULT 0,
    position_y FLOAT NOT NULL DEFAULT 0,
    color CHAR(7) NOT NULL DEFAULT '#FF0000',
    texture_path VARCHAR(255),
    experience_points INTEGER NOT NULL DEFAULT 0,
    clan_id INTEGER NULL DEFAULT NULL,
    PRIMARY KEY(id),
    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(clan_id) REFERENCES clans(id)
);
CREATE TABLE weapon_classes (
    id SERIAL,
    base_damage INTEGER NOT NULL,
    base_range FLOAT NOT NULL,
    base_rate_of_fire INTEGER NOT NULL,
    PRIMARY KEY(id)
);
CREATE TABLE weapons (
    id SERIAL,
    weapon_class_id INTEGER NOT NULL,
    level INTEGER NOT NULL DEFAULT 1,
    player_id INTEGER NOT NULL,
    PRIMARY KEY(id),
    FOREIGN KEY(weapon_class_id) REFERENCES weapon_classes(id),
    FOREIGN KEY(player_id) REFERENCES players(id)
);
INSERT INTO weapon_classes (base_damage, base_range, base_rate_of_fire)
VALUES (75, 1.0, 6),
    -- Katana: one hit every 75 ticks
    (25, 10.0, 2),
    -- Pistol: one hit every 45 ticks
    (90, 7.5, 5),
    -- Shotgun: one hit every 45 ticks
    (55, 20.0, 3),
    -- Rifle: one hit every 20 ticks
    (120, 15.0, 10);
-- Hand Grenade: one hit every 75 ticks