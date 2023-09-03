package db

import (
	"context"
	"fmt"
	"wut/ent"
	"wut/ent/user"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type DbClient struct {
	Client *ent.Client
}

func SetupDBConnection() *DbClient {

	client, err := ent.Open("postgres", "host=localhost port=5432 user=myuser dbname=test1 password=mypassword sslmode=disable")
	if err != nil {
		fmt.Println("failed opening connection to postgres: %w", err)
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		fmt.Println("failed creating schema resources: %w", err)
	}
	return &DbClient{client}
}

func (dbI *DbClient) CreateUser(ctx context.Context, email, password string) (*ent.User, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed hashing password: %w", err)
	}
	fmt.Println(hash, "  ", password)
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

func (dbI *DbClient) AuthUser(ctx context.Context, email, password string) (*ent.User, error) {
	user, err := dbI.Client.User.Query().
		Where(user.EmailEQ(email)).
		Only(ctx)
	if ent.IsNotFound(err) {
		return nil, fmt.Errorf("email doesn't exist")
	} else if err != nil {
		return nil, fmt.Errorf("server side error")
	}
	err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password))
	if err == nil {
		return user, nil
	} else {
		return nil, fmt.Errorf("password does not match")
	}
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
