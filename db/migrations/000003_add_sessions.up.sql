CREATE TABLE sessions (
    id uuid PRIMARY KEY,
    username varchar(255) NOT NULL,
    refresh_token varchar(255) NOT NULL,
    user_agent varchar(255) NOT NULL,
    client_ip varchar(255) NOT NULL,
    id_blocked boolean NOT NULL DEFAULT false,
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE sessions ADD FOREIGN KEY (username) REFERENCES users (username);