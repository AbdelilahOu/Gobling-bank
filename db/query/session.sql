-- name: CreateSession :one
insert into sessions (
  id,
  username,
  refresh_token,
  user_agent,
  client_ip,
  id_blocked,
  expires_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions
WHERE username = $1 LIMIT 1;