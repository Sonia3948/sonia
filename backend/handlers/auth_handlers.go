package handlers

import (
	"elsaidaliya/handlers/auth"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// InitAuthHandlers initializes the collections for auth handlers
func InitAuthHandlers(db *mongo.Database) {
	auth.Init(db)
}

// RegisterUser handles user registration
func RegisterUser(c *gin.Context) {
	auth.RegisterUser(c)
}

// LoginUser function handles both regular user and admin login
func LoginUser(c *gin.Context) {
	auth.LoginUser(c)
}

// ForgotPassword handles password reset requests
func ForgotPassword(c *gin.Context) {
	auth.ForgotPassword(c)
}
