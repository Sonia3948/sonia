
package offer

import (
	"context"
	"net/http"
	"time"

	"elsaidaliya/middleware"
	"elsaidaliya/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Update updates an existing offer
func Update(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	offerID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID d'offre invalide"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if offer exists and belongs to user
	var existingOffer models.Offer
	err = offerCollection.FindOne(ctx, bson.M{
		"_id":        offerID,
		"supplierID": userID,
	}).Decode(&existingOffer)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Offre non trouvée ou vous n'avez pas l'autorisation de la modifier"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la vérification de l'offre"})
		return
	}

	// Bind the update data
	var updateData models.Offer
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	// Create update document
	update := bson.M{
		"$set": bson.M{
			"title":       updateData.Title,
			"description": updateData.Description,
			"price":       updateData.Price,
			"imageUrl":    updateData.ImageURL,
			"expiresAt":   updateData.ExpiresAt,
			"updatedAt":   time.Now(),
		},
	}

	// Update the offer
	_, err = offerCollection.UpdateOne(ctx, bson.M{"_id": offerID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour de l'offre"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Offre mise à jour avec succès"})
}
