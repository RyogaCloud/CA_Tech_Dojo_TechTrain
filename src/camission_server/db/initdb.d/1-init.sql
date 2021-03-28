CREATE DATABASE IF NOT EXISTS ca_mission;

USE ca_mission;

CREATE TABLE IF NOT EXISTS users (
    id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    token VARCHAR(100) NOT NULL,
    created_at timestamp NOT NULL default current_timestamp,
    updated_at timestamp NOT NULL default current_timestamp on update current_timestamp
);

CREATE TABLE IF NOT EXISTS charas (
    id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    probability int NOT NULL
);

CREATE TABLE IF NOT EXISTS users_charas_possession (
    id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id int NOT NULL,
    chara_id int NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    FOREIGN KEY (chara_id) REFERENCES charas(id) ON DELETE RESTRICT ON UPDATE CASCADE
);

INSERT INTO
    users (name, token)
VALUES
    ("test1", "aaaa"),
    ("test2", "bbbb"),
    ("test3", "cccc");

INSERT INTO
  charas (name, probability)
VALUES
  ("めちゃめちゃレア1", 1),
  ("めちゃめちゃレア２", 1),
  ("めちゃレア１", 2),
  ("めちゃレア2", 2),
  ("めちゃレア3", 2),
  ("めちゃレア4", 2),
  ("レア1", 3),
  ("レア2", 3),
  ("レア3", 3),
  ("レア4", 3),
  ("レア5", 3),
  ("レア6", 3),
  ("レア7", 3),
  ("レア8", 3);