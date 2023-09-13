CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    pass_hash BYTEA NULL,
    created_date DATE NOT NULL DEFAULT CURRENT_DATE,
    last_jwt_revoke_date DATE NOT NULL DEFAULT CURRENT_DATE,
    username varchar(20) NOT NULL,
    order_num varchar(36) NOT NULL,
    email varchar(254) UNIQUE NOT NULL
);

-- atlas schema apply -u "postgres://myuser@localhost:5432/test1?search_path=public&sslmode=disable" --to "file://./schema.sql" --dev-url "postgres://postgres@localhost:5433/postgres?search_path=public&sslmode=disable"
