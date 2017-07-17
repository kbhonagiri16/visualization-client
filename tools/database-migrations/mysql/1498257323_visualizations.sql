-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE visualization (
    id int unsigned NOT NULL AUTO_INCREMENT,
    name Varchar(255) NOT NULL,
    slug Varchar(36) NOT NULL,
    organization_id Varchar(36) NOT NULL,
    tags json DEFAULT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE dashboard (
    id Varchar(36) NOT NULL,
    visualization_id int unsigned NOT NULL,
    name Varchar(255) NOT NULL,
    slug Varchar(255) NOT NULL,
    rendered_template MEDIUMTEXT NOT NULL,
    PRIMARY KEY(id),
    FOREIGN KEY (visualization_id)
        REFERENCES visualization(id)
        ON DELETE CASCADE
);


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE dashboard;
DROP TABLE visualization;
