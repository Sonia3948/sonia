
package routes

import (
	"elsaidaliya/handlers"

	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes configures routes for authentication
func SetupAuthRoutes(r *gin.Engine) {
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", handlers.RegisterUser)
		auth.POST("/login", handlers.LoginUser)
		auth.POST("/forgot-password", handlers.ForgotPassword)
	}
}
