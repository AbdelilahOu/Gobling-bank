CREATE TABLE users (
    username varchar(255) UNIQUE NOT NULL PRIMARY KEY,
    hashed_password varchar(255) NOT NULL,
    full_name varchar(255) NOT NULL,
    email varchar(255) UNIQUE NOT NULL,
    password_changed_at timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    created_at timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE accounts ADD FOREIGN KEY (owner) REFERENCES users (username);

ALTER TABLE accounts ADD CONSTRAINT owner_currency_key UNIQUE (owner,currency);