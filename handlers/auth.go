package handlers

import (
	"Go-Password-Manager/models"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{DB: db}
}

func (h *AuthHandler) ShowLoginPage(c *gin.Context) {
	// If already logged in, redirect to dashboard
	session := sessions.Default(c)
	if userID := session.Get("userID"); userID != nil {
		c.Redirect(http.StatusFound, "/dashboard")
		return
	}

	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func (h *AuthHandler) ProcessLogin(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	var user models.User
	if err := h.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.MasterPassword), []byte(password)); err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Set session with user ID
	session := sessions.Default(c)
	session.Set("userID", user.ID)
	if err := session.Save(); err != nil {
		log.Println("Error saving session:", err)
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{
			"error": "Unable to create session",
		})
		return
	}

	c.Redirect(http.StatusFound, "/dashboard")
}

func (h *AuthHandler) ShowRegisterPage(c *gin.Context) {
	// If already logged in, redirect to dashboard
	session := sessions.Default(c)
	if userID := session.Get("userID"); userID != nil {
		c.Redirect(http.StatusFound, "/dashboard")
		return
	}

	c.HTML(http.StatusOK, "register.html", gin.H{})
}

func (h *AuthHandler) ProcessRegister(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	confirmPassword := c.PostForm("confirm_password")

	// Validation
	if email == "" || password == "" || confirmPassword == "" {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error": "All fields are required",
			"email": email,
		})
		return
	}

	if password != confirmPassword {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error": "Passwords do not match",
			"email": email,
		})
		return
	}

	if len(password) < 8 {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error": "Password must be at least 8 characters long",
			"email": email,
		})
		return
	}

	// Check if user already exists
	var existingUser models.User
	if err := h.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error": "An account with this email already exists",
			"email": email,
		})
		return
	}

	// Hash the master password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{
			"error": "Unable to create account. Please try again.",
			"email": email,
		})
		return
	}

	// Create new user
	user := models.User{
		Email:          email,
		MasterPassword: string(hashedPassword),
	}

	if err := h.DB.Create(&user).Error; err != nil {
		log.Println("Error creating user:", err)
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{
			"error": "Unable to create account. Please try again.",
			"email": email,
		})
		return
	}

	log.Printf("New user registered: %s (ID: %d)", user.Email, user.ID)

	// Auto-login after registration
	session := sessions.Default(c)
	session.Set("userID", user.ID)
	if err := session.Save(); err != nil {
		log.Println("Error saving session:", err)
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.Redirect(http.StatusFound, "/dashboard")
}

func (h *AuthHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		log.Println("Unable to clear session:", err)
	}

	c.Redirect(http.StatusFound, "/login")
}
