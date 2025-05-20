package handlers

import (
	"context"
	"net/http"
	"time"

	"elsaidaliya/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection *mongo.Collection

// InitUserHandlers initializes handlers for user management
func InitUserHandlers(db *mongo.Database) {
	userCollection = db.Collection("users")
}

// GetAllUsers fetches all users
func GetAllUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Define options to exclude password field
	findOptions := options.Find()
	findOptions.SetProjection(bson.M{"password": 0})

	// Get filter from query params
	filter := bson.M{}
	if role := c.Query("role"); role != "" {
		filter["role"] = role
	}
	if active := c.Query("isActive"); active != "" {
		if active == "true" {
			filter["isActive"] = true
		} else {
			filter["isActive"] = false
		}
	}
	if subscription := c.Query("subscription"); subscription != "" {
		filter["subscription"] = subscription
	}

	cursor, err := userCollection.Find(ctx, filter, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des utilisateurs"})
		return
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du traitement des utilisateurs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// GetFeaturedSuppliers fetches suppliers with gold subscription
func GetFeaturedSuppliers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Define options to exclude sensitive fields
	findOptions := options.Find()
	findOptions.SetProjection(bson.M{
		"password": 0,
		"subscription": 0,  // Hide subscription info
		"subExpiry": 0,     // Hide subscription expiry
	})

	// Filter for active suppliers with gold subscription
	filter := bson.M{
		"role":         "fournisseur",
		"isActive":     true,
		"subscription": "or",  // Only fetch gold subscription suppliers
	}

	cursor, err := userCollection.Find(ctx, filter, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des fournisseurs vedettes"})
		return
	}
	defer cursor.Close(ctx)

	var suppliers []models.User
	if err = cursor.All(ctx, &suppliers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du traitement des fournisseurs vedettes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"suppliers": suppliers})
}

// GetPendingUsers fetches users pending approval
func GetPendingUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find users where isActive is false (pending approval)
	filter := bson.M{"isActive": false}
	
	// Define options to exclude password field
	findOptions := options.Find()
	findOptions.SetProjection(bson.M{"password": 0})
	
	cursor, err := userCollection.Find(ctx, filter, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des utilisateurs en attente"})
		return
	}
	defer cursor.Close(ctx)

	var pendingUsers []models.User
	if err = cursor.All(ctx, &pendingUsers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du traitement des utilisateurs en attente"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": pendingUsers})
}

// GetUserByID fetches a single user by ID
func GetUserByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID utilisateur invalide"})
		return
	}

	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération de l'utilisateur"})
		return
	}

	// Don't return the password
	user.Password = ""

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateUserStatus updates a user's active status
func UpdateUserStatus(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID utilisateur invalide"})
		return
	}

	var statusUpdate struct {
		IsActive bool `json:"isActive" binding:"required"`
	}

	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"isActive":  statusUpdate.IsActive,
			"updatedAt": time.Now(),
		},
	}

	result, err := userCollection.UpdateOne(ctx, bson.M{"_id": userID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour du statut"})
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Statut utilisateur mis à jour avec succès"})
}

// UpdateUserSubscription updates a user's subscription details
func UpdateUserSubscription(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID utilisateur invalide"})
		return
	}

	var subUpdate struct {
		Subscription string    `json:"subscription" binding:"required"`
		SubExpiry    time.Time `json:"subExpiry" binding:"required"`
	}

	if err := c.ShouldBindJSON(&subUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"subscription": subUpdate.Subscription,
			"subExpiry":    subUpdate.SubExpiry,
			"updatedAt":    time.Now(),
		},
	}

	result, err := userCollection.UpdateOne(ctx, bson.M{"_id": userID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour de l'abonnement"})
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Abonnement utilisateur mis à jour avec succès"})
}

// UpdateUser updates user details
func UpdateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID utilisateur invalide"})
		return
	}

	var update struct {
		BusinessName   string `json:"businessName"`
		Phone          string `json:"phone"`
		Email          string `json:"email"`
		Wilaya         string `json:"wilaya"`
		IsActive       bool   `json:"isActive"`
		RegisterNumber string `json:"registerNumber"`
		Subscription   string `json:"subscription"`
		SubExpiry      string `json:"subExpiry"`
	}

	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	updateFields := bson.M{
		"updatedAt": time.Now(),
	}

	if update.BusinessName != "" {
		updateFields["businessName"] = update.BusinessName
	}
	if update.Phone != "" {
		updateFields["phone"] = update.Phone
	}
	if update.Email != "" {
		updateFields["email"] = update.Email
	}
	if update.Wilaya != "" {
		updateFields["wilaya"] = update.Wilaya
	}
	if update.RegisterNumber != "" {
		updateFields["registerNumber"] = update.RegisterNumber
	}
	if update.Subscription != "" {
		updateFields["subscription"] = update.Subscription
	}
	if update.SubExpiry != "" {
		// Parse the date string to time.Time
		expiry, err := time.Parse("2006-01-02", update.SubExpiry)
		if err == nil {
			updateFields["subExpiry"] = expiry
		}
	}
	
	updateFields["isActive"] = update.IsActive

	updateDoc := bson.M{
		"$set": updateFields,
	}

	result, err := userCollection.UpdateOne(ctx, bson.M{"_id": userID}, updateDoc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour de l'utilisateur"})
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé ou aucun changement appliqué"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Utilisateur mis à jour avec succès"})
}

// GetUserRegisterImage returns the register image URL for a specific user
func GetUserRegisterImage(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID utilisateur invalide"})
		return
	}

	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération de l'utilisateur"})
		return
	}

	if user.RegisterImageURL == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image de registre de commerce non trouvée"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"imageUrl": user.RegisterImageURL})
}
