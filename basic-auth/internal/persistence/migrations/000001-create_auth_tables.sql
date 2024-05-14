-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE auth_users (
    id integer primary key autoincrement,
    user varchar(50) not null unique,
    pass varchar(255) not null
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE auth_users;
