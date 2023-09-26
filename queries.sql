-- name: CreateUser :exec
INSERT INTO users (email, username, payment_intent, donation_in_cents ) VALUES ($1, $2, $3, $4);

-- name: GetUserPassAndId :one
SELECT id, pass_hash FROM users where email = $1 LIMIT 1;

-- name: UpdateUserPass :exec
UPDATE users SET pass_hash = $2 WHERE email = $1;

-- name: InitPass :one
UPDATE users SET pass_hash=$2 WHERE payment_intent=$1 AND pass_hash IS NULL RETURNING id;

-- name: CountUsers :one
SELECT COUNT(id) FROM users;

-- name: GetNUsers :many
SELECT ROW_NUMBER() OVER (ORDER BY id), username, donation_in_cents, created_date FROM users ORDER BY id::int ASC LIMIT @n OFFSET (@page - 1) * @n::int;

-- name: GetUserPos :one
SELECT COUNT(id) FROM users WHERE id<=$1;

-- name: AddToBalance :exec
UPDATE users SET donation_in_cents=donation_in_cents + @add_amount::int WHERE email = @email::varchar(254);