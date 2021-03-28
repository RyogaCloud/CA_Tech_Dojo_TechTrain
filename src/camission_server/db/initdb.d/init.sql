CREATE DATABASE IF NOT EXISTS ca_mission;

USE ca_mission;

CREATE TABLE IF NOT EXISTS user (
    id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    token VARCHAR(100) NOT NULL,
    created_at timestamp NOT NULL default current_timestamp,
    updated_at timestamp NOT NULL default current_timestamp on update current_timestamp,
);

INSERT INTO
    user (name, token)
VALUES
    ("test1", "aaaa"),
    ("test2", "bbbb"),
    ("test3", "cccc");