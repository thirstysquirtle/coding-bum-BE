package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	fmt.Println(bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost))

}
