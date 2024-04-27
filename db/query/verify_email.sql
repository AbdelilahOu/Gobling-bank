-- name: CreateVerifyEmail :one
INSERT INTO verify_emails (username, email, secret_code) VALUES ($1, $2, $3) RETURNING *;

-- name: GetVerifyEmail :one
SELECT * FROM verify_emails WHERE id = $1;

-- name: UpdateVerifyEmail :one
UPDATE verify_emails SET is_used = TRUE WHERE id = @id AND secret_code = @secret_code AND is_used = FALSE AND expired_at > now() RETURNING *;