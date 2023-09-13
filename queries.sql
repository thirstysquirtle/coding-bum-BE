-- name: CreateUser :exec
INSERT INTO users (email, username, payment_intent ) VALUES ($1, $2, $3);

-- name: GetUserPass :one
SELECT pass_hash FROM users where email = $1 LIMIT 1;

-- name: UpdateUserPass :exec
UPDATE users SET pass_hash = $2 WHERE email = $1;
