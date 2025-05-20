
package listing

import (
	"context"
	"net/http"
	"time"

	"elsaidaliya/middleware"
	"elsaidaliya/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Create creates a new medicine listing
func Create(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	var newListing models.Listing
	if err := c.ShouldBindJSON(&newListing); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides", "details": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set the supplier ID to the current user
	newListing.ID = primitive.NewObjectID()
	newListing.SupplierID = userID
	newListing.CreatedAt = time.Now()
	newListing.UpdatedAt = time.Now()

	// Insert listing
	_, err = listingCollection.InsertOne(ctx, newListing)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création du listing"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Listing créé avec succès",
		"listing": newListing,
	})
}
