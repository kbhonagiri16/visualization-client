-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE user (
    id int NOT NULL AUTO_INCREMENT,
    usr_name Varchar(255) NOT NULL,
    PRIMARY KEY(id)
);


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE user;
