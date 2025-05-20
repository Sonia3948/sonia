package routes

import (
	"elsaidaliya/handlers"
	"elsaidaliya/middleware"

	"github.com/gin-gonic/gin"
)

// SetupNotificationRoutes initializes routes for notification management
func SetupNotificationRoutes(r *gin.Engine) {
	notificationGroup := r.Group("/api/notifications")
	notificationGroup.POST("/payment", handlers.CreatePaymentNotification)
	notificationGroup.POST("/create", middleware.RequireAuth, handlers.CreateNotification)
	
	// Protected routes requiring authentication
	notificationGroup.Use(middleware.RequireAuth)
	{
		notificationGroup.GET("/admin", handlers.GetAdminNotifications)
		notificationGroup.GET("/user", handlers.GetUserNotifications)
		notificationGroup.PUT("/:id/read", handlers.MarkNotificationAsRead)
		notificationGroup.PUT("/:id/status", handlers.UpdateNotificationStatus)
	}
}
