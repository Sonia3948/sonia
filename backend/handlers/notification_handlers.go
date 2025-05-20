package handlers

import (
	"elsaidaliya/handlers/notification"

	"go.mongodb.org/mongo-driver/mongo"
)

// InitNotificationHandlers initializes handlers for notification management
func InitNotificationHandlers(db *mongo.Database) {
	notification.InitNotificationHandlers(db)
}

// CreateNotification creates a new notification
var CreateNotification = notification.CreateNotification

// CreatePaymentNotification creates a notification when a payment receipt is uploaded
var CreatePaymentNotification = notification.CreatePaymentNotification

// GetUserNotifications gets all notifications for a user
var GetUserNotifications = notification.GetUserNotifications

// GetAdminNotifications gets all notifications for admin
var GetAdminNotifications = notification.GetAdminNotifications

// MarkNotificationAsRead marks a notification as read
var MarkNotificationAsRead = notification.MarkNotificationAsRead

// UpdateNotificationStatus updates the status of a notification
var UpdateNotificationStatus = notification.UpdateNotificationStatus
