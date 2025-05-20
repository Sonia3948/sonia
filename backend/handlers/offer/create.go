
package offer

import (
	"context"
	"net/http"
	"time"

	"elsaidaliya/middleware"
	"elsaidaliya/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Create creates a new offer
func Create(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	var newOffer models.Offer
	if err := c.ShouldBindJSON(&newOffer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides", "details": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set the supplier ID to the current user
	newOffer.ID = primitive.NewObjectID()
	newOffer.SupplierID = userID
	newOffer.CreatedAt = time.Now()
	newOffer.UpdatedAt = time.Now()

	// Insert offer
	_, err = offerCollection.InsertOne(ctx, newOffer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de l'offre"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Offre créée avec succès",
		"offer":   newOffer,
	})
}
