package auth

import (
	"context"
	"log"
	"net/http"
	"time"

	"elsaidaliya/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser handles user registration
func RegisterUser(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides", "details": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if user already exists
	var existingUser models.User
	err := userCollection.FindOne(ctx, bson.M{
		"$or": []bson.M{
			{"email": newUser.Email},
			{"phone": newUser.Phone},
		},
	}).Decode(&existingUser)

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Un utilisateur avec cet e-mail ou ce numéro de téléphone existe déjà"})
		return
	} else if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la vérification de l'utilisateur"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du traitement du mot de passe"})
		return
	}

	// Prepare user document
	newUser.ID = primitive.NewObjectID()
	newUser.Password = string(hashedPassword)
	newUser.IsActive = false // User requires activation
	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()

	if newUser.Role == "" {
		newUser.Role = "pharmacien" // Default role
	}

	// Insert user
	_, err = userCollection.InsertOne(ctx, newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de l'utilisateur"})
		return
	}

	// Create notification for admin about new user registration
	adminNotification := models.Notification{
		ID:          primitive.NewObjectID(),
		UserID:      primitive.ObjectID{}, // Will be set to admin ID in production
		FromID:      newUser.ID,
		FromName:    newUser.BusinessName,
		Type:        "user_registration",
		Title:       "Nouvelle inscription",
		Description: "Un nouveau " + newUser.Role + " s'est inscrit et attend validation",
		Read:        false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save notification to database
	_, err = userCollection.Database().Collection("notifications").InsertOne(ctx, adminNotification)
	if err != nil {
		// Log error but don't fail the registration
		log.Printf("Error creating admin notification: %v", err)
	}

	// Don't return the password
	newUser.Password = ""

	c.JSON(http.StatusCreated, gin.H{
		"message": "Utilisateur enregistré avec succès. Votre compte est en attente d'activation par un administrateur.",
		"user":    newUser,
	})
}
