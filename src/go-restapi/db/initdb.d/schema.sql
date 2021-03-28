CREATE TABLE Comment
(
    id int NOT NULL AUTO_INCREMENT,
    title varchar(32) NOT NULL,
    category varchar(32) NOT NULL,
    content varchar(256) NOT NULL,
    PRIMARY KEY(id)
);