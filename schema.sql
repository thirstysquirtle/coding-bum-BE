CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    donation_in_cents INT NOT NULL,
    created_date DATE NOT NULL DEFAULT CURRENT_DATE,
    last_jwt_revoke_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    username varchar(20) NOT NULL,
    payment_intent varchar(36) NOT NULL UNIQUE,
    email varchar(254) UNIQUE NOT NULL,
    pass_hash BYTEA NULL

);

-- atlas schema apply -u "postgres://myuser@localhost:5434/test1?search_path=public&sslmode=disable" --to "file://./schema.sql" --dev-url "postgres://postgres@localhost:5433/postgres?search_path=public&sslmode=disable"
