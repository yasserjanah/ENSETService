package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	ID, err := uuid.NewV4()
	if err != nil {
		log.Fatalln("ERROR :", err.Error())
	}
	var Email string = "admin@test.ma"
	var Password string = "password_test"
	fmt.Println("*** Generating bcrypt hash...")
	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalln("ERROR :", err.Error())
	}
	createdAt := time.Now().Format(time.RFC3339)
	updatedAt := time.Now().Format(time.RFC3339)

	fmt.Println("*** Generating SQL statement...")
	fmt.Println()
	fmt.Printf("INSERT INTO admins(id, email, password_hash, password, password_confirmation, created_at, updated_at) VALUES ('%v', '%v', '%v', '-', '-', '%v', '%v');\n", ID.String(), Email, string(HashedPassword), createdAt, updatedAt)
	fmt.Println()
}
