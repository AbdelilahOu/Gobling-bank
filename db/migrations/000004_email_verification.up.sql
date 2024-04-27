CREATE TABLE verify_emails (
    id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    username varchar(255) NOT NULL,
    email varchar(255) NOT NULL,
    secret_code varchar(255) NOT NULL,
    is_used bool NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT (now()),
    expired_at timestamptz NOT NULL DEFAULT (now() + interval '15 minutes')
);

ALTER TABLE verify_emails ADD FOREIGN KEY (username) REFERENCES users (username);

ALTER TABLE users ADD COLUMN is_email_verified bool NOT NULL DEFAULT false;

