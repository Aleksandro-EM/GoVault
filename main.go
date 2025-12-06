package main

import (
	"Go-Password-Manager/database"
	"Go-Password-Manager/handlers"
	"Go-Password-Manager/middleware"
	"Go-Password-Manager/models"
	"log"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	config := database.LoadConfig()

	db, err := database.Connect(config)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Password{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database connected and migrated!")

	router := gin.Default()

	// Validate and setup sessions
	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatal("SECRET environment variable is required")
	}
	if len(secret) < 32 {
		log.Fatal("SECRET must be at least 32 characters")
	}

	store := cookie.NewStore([]byte(secret))
	store.Options(sessions.Options{
		MaxAge:   60 * 60 * 24,
		HttpOnly: true,
		Secure:   false, //set to true in prod
		Path:     "/",
	})
	router.Use(sessions.Sessions("govault", store))

	router.LoadHTMLGlob("templates/*")

	authHandler := handlers.NewAuthHandler(db)

	// Public routes
	router.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/login")
	})
	router.GET("/login", authHandler.ShowLoginPage)
	router.POST("/login", authHandler.ProcessLogin)
	router.GET("/register", authHandler.ShowRegisterPage)
	router.POST("/register", authHandler.ProcessRegister)

	// Protected routes
	protected := router.Group("/")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/dashboard", func(c *gin.Context) {
			session := sessions.Default(c)
			userID := session.Get("userID")

			c.HTML(200, "dashboard.html", gin.H{
				"message": "Welcome to GoVault!",
				"userID":  userID,
			})
		})
		protected.GET("/logout", authHandler.Logout)
	}

	log.Println("Server starting on :8080")
	router.Run(":8080")
}
