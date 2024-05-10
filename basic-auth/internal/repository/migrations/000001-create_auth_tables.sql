-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE auth_users (
    id integer primary key autoincrement,
    user varchar(50) not null unique,
    pass varchar(255) not null
);

insert into auth_users ("user", pass) values (
    'root'
    , '$2a$10$Aq3G94WIe2JBQAXQG3G4V.oFd.k6Noe3MAXcp2VCTLJx/.j09LZqa' -- 12345
);

insert into auth_users ("user", pass) values (
    'user'
    , '$2a$10$Aq3G94WIe2JBQAXQG3G4V.oFd.k6Noe3MAXcp2VCTLJx/.j09LZqa' -- 12345
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE auth_users;