CREATE DATABASE IF NOT EXISTS ca_mission;

USE ca_mission;

CREATE TABLE IF NOT EXISTS users (
    id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    token VARCHAR(100) NOT NULL,
    created_at timestamp NOT NULL default current_timestamp,
    updated_at timestamp NOT NULL default current_timestamp on update current_timestamp,
    deleted_at timestamp default NULL
);

CREATE TABLE IF NOT EXISTS characters (
    id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    probability int NOT NULL,
    created_at timestamp NOT NULL default current_timestamp,
    updated_at timestamp NOT NULL default current_timestamp on update current_timestamp,
    deleted_at timestamp default NULL
);

CREATE TABLE IF NOT EXISTS user_has_characters (
    id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id int NOT NULL,
    character_id int NOT NULL,
    created_at timestamp NOT NULL default current_timestamp,
    updated_at timestamp NOT NULL default current_timestamp on update current_timestamp,
    deleted_at timestamp default NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE RESTRICT ON UPDATE CASCADE
);

INSERT INTO
    users (name, token)
VALUES
    ('testuser1', 'aaaa'),
    ('testuser2', 'bbbb'),
    ('testuser3', 'cccc'),
    ('testuser4', 'dddd');

INSERT INTO
    chara (name, probability)
VALUES
    ("UR1", 1),
    ("UR２", 1),
    ("SR１", 3),
    ("SR2", 3),
    ("SR3", 3),
    ("SR4", 3),
    ("R1", 6),
    ("R2", 6),
    ("R3", 6),
    ("R4", 6),
    ("R5", 6),
    ("R6", 6),
    ("R7", 6),
    ("R8", 6);