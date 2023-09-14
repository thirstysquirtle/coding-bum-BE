-- name: CreateUser :exec
INSERT INTO users (email, username, payment_intent, donation_in_cents ) VALUES ($1, $2, $3, $4);

-- name: GetUserPass :one
SELECT pass_hash FROM users where email = $1 LIMIT 1;

-- name: UpdateUserPass :exec
UPDATE users SET pass_hash = $2 WHERE email = $1;

-- name: InitPass :one
UPDATE users SET pass_hash=$2 WHERE payment_intent=$1 AND pass_hash IS NULL RETURNING id;