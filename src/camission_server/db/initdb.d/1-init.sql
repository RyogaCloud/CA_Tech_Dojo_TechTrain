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

INSERT INTO
    users (name, token)
VALUES
    ('testuser1', '1616916651'),
    ('testuser2', '1516916651'),
    ('testuser3', '1416916651'),
    ('testuser4', '1316916651');