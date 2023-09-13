-- name: CreateUser :copyfrom
INSERT INTO users (email, username, order_num ) VALUES ($1, $2, $3);

-- name: getUserPass :one
SELECT pass_hash FROM users where email = $1 LIMIT 1;

-- name: updateUserPass :exec
UPDATE users SET pass_hash = $2 WHERE email = $1;
