package notification

import (
	"context"
	"elsaidaliya/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Collection holds the MongoDB collection for notifications
var Collection *mongo.Collection

// InitNotificationHandlers initializes the notification handlers
func InitNotificationHandlers(db *mongo.Database) {
	Collection = db.Collection("notifications")
}

// CreateNotification creates a new notification
func CreateNotification(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var notif models.Notification
	if err := c.ShouldBindJSON(&notif); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notif.CreatedAt = time.Now()
	notif.UpdatedAt = time.Now()
	notif.Read = false

	result, err := Collection.InsertOne(ctx, notif)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de la notification"})
		return
	}

	notif.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, gin.H{"notification": notif})
}

// CreatePaymentNotification creates a payment notification
func CreatePaymentNotification(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var data struct {
		UserID      string `json:"userID"`
		UserName    string `json:"userName"`
		PaymentID   string `json:"paymentID"`
		Amount      string `json:"amount"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := primitive.ObjectIDFromHex(data.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UserID invalide"})
		return
	}

	adminID := primitive.NewObjectID() // Placeholder, notification for admin

	notification := models.Notification{
		UserID:      adminID,
		FromID:      userID,
		FromName:    data.UserName,
		Type:        models.NotificationTypePayment,
		Title:       "Nouveau bon de versement",
		Description: "Un bon de versement de " + data.Amount + " DZD a été téléversé",
		Read:        false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result, err := Collection.InsertOne(ctx, notification)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de la notification"})
		return
	}

	notification.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, gin.H{"notification": notification})
}

// GetUserNotifications retrieves notifications for a user
func GetUserNotifications(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID := c.MustGet("userID").(primitive.ObjectID)

	// Define options to sort by createdAt DESC
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"createdAt": -1})

	cursor, err := Collection.Find(ctx, bson.M{"userID": userID}, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des notifications"})
		return
	}
	defer cursor.Close(ctx)

	var notifications []models.Notification
	if err = cursor.All(ctx, &notifications); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du traitement des notifications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notifications": notifications})
}

// GetAdminNotifications retrieves notifications for admin
func GetAdminNotifications(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get admin ID from token
	adminID := c.MustGet("userID").(primitive.ObjectID)

	// Verify user is admin (should be handled by middleware)
	role := c.MustGet("userRole").(string)
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Accès non autorisé"})
		return
	}

	// Define options to sort by createdAt DESC
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"createdAt": -1})

	cursor, err := Collection.Find(ctx, bson.M{"userID": adminID}, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des notifications"})
		return
	}
	defer cursor.Close(ctx)

	var notifications []models.Notification
	if err = cursor.All(ctx, &notifications); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du traitement des notifications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notifications": notifications})
}

// MarkNotificationAsRead marks a notification as read
func MarkNotificationAsRead(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	notifID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de notification invalide"})
		return
	}

	userID := c.MustGet("userID").(primitive.ObjectID)

	filter := bson.M{"_id": notifID, "userID": userID}
	update := bson.M{"$set": bson.M{"read": true, "updatedAt": time.Now()}}

	result, err := Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour de la notification"})
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification non trouvée ou déjà marquée comme lue"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marquée comme lue"})
}

// UpdateNotificationStatus updates the status of a notification
func UpdateNotificationStatus(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	notifID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de notification invalide"})
		return
	}

	var data struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(primitive.ObjectID)

	filter := bson.M{"_id": notifID, "userID": userID}
	update := bson.M{"$set": bson.M{"status": data.Status, "updatedAt": time.Now()}}

	result, err := Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour du statut"})
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification non trouvée"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Statut de la notification mis à jour"})
}
