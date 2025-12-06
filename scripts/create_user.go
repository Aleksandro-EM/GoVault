package main

import (
	"Go-Password-Manager/database"
	"Go-Password-Manager/models"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	config := database.LoadConfig()
	db, err := database.Connect(config)
	if err != nil {
		log.Fatal(err)
	}

	// Hash password
	hash, _ := bcrypt.GenerateFromPassword([]byte("test123"), bcrypt.DefaultCost)

	user := models.User{
		Email:          "test@example.com",
		MasterPassword: string(hash),
	}

	db.Create(&user)
	fmt.Println("Test user created: test@example.com / test123")
}
