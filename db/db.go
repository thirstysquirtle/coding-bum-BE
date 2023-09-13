package db

import (
	"context"
	"fmt"
	"os"
	"wut/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

var PostgresPool *pgxpool.Pool
var Db *sqlc.Queries

func SetupDBConnection() {
	setEnv()
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(fmt.Sprintf("Unable to create connection pool to PostgreSQL: %v\n", err))
	}
	PostgresPool = dbpool
	Db = sqlc.New(PostgresPool)
}

// func (dbI *DbClient) CreateUser(ctx context.Context, email, password string) (*ent.User, error) {

// 	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed hashing password: %w", err)
// 	}
// 	fmt.Println(hash, "  ", password)
// 	u, err := dbI.Client.User.
// 		Create().
// 		SetUsername("test").
// 		SetPassHash(hash).
// 		SetEmail(email).
// 		SetOrderID("test").
// 		Save(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed saving user: %w", err)
// 	}
// 	// fmt.Println("user was created: ", u)
// 	return u, nil
// }

// func (dbI *DbClient) AuthUser(ctx context.Context, email, password string) (*ent.User, error) {
// 	user, err := dbI.Client.User.Query().
// 		Where(user.EmailEQ(email)).
// 		Only(ctx)
// 	if ent.IsNotFound(err) {
// 		return nil, fmt.Errorf("email doesn't exist")
// 	} else if err != nil {
// 		return nil, fmt.Errorf("server side error")
// 	}
// 	err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password))
// 	if err == nil {
// 		return user, nil
// 	} else {
// 		return nil, fmt.Errorf("password does not match")
// 	}
// }

func setEnv() {
	if os.Getenv("IsTest") == "true" {
		os.Setenv("DATABASE_URL", "postgres://postgres@localhost:5434/test1?sslmode=disable")

	}
}
