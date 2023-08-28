package db

import (
	"context"
	"fmt"
	"wut/ent"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type DbClient struct {
	Client *ent.Client
}

func SetupDBConnection() *DbClient {

	client, err := ent.Open("postgres", "host=localhost port=5432 user=myuser dbname=test1 password=mypassword sslmode=disable")
	if err != nil {
		fmt.Printf("failed opening connection to postgres: %v", err)
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		fmt.Printf("failed creating schema resources: %v", err)
	}
	return &DbClient{client}
}

func (dbI *DbClient) CreateUser(ctx context.Context, email, password string) (*ent.User, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed hashing password: %w", err)
	}

	u, err := dbI.Client.User.
		Create().
		SetUsername("test").
		SetPassHash(hash).
		SetEmail(email).
		SetOrderID("test").
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed saving user: %w", err)
	}
	// fmt.Println("user was created: ", u)
	return u, nil
}

// func getSecret() []byte {
// 	if os.Getenv("IsTest") == "true" {
// 		os.Setenv("secret1", "testSecret")
// 	}

// 	secretKeyStr, found := os.LookupEnv("secret1")
// 	if !found {
// 		panic("No secret Env Variable Found!")
// 	}
// 	return []byte(secretKeyStr)
// }
