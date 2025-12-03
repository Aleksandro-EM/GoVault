package main

import (
	"Go-Password-Manager/database"
	"Go-Password-Manager/models"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load config and connect to DB
	config := database.LoadConfig()

	db, err := database.Connect(config)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate models (creates/updates tables)
	err = db.AutoMigrate(&models.User{}, &models.Password{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database connected and migrated!")

	// Initialize Gin router
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":   "ok",
			"database": "connected",
		})
	})

	log.Println("Server starting on :8080")
	router.Run(":8080")
}
