
package listing

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

// Update updates an existing medicine listing
func Update(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	listingID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de listing invalide"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if listing exists and belongs to user
	var existingListing models.Listing
	err = listingCollection.FindOne(ctx, bson.M{
		"_id": listingID,
		"supplierID": userID,
	}).Decode(&existingListing)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Listing non trouvé ou vous n'avez pas l'autorisation de le modifier"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la vérification du listing"})
		return
	}

	// Bind the update data
	var updateData models.Listing
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	// Create update document
	update := bson.M{
		"$set": bson.M{
			"title":       updateData.Title,
			"description": updateData.Description,
			"medications": updateData.Medications,
			"pdfUrl":      updateData.PdfURL,
			"updatedAt":   time.Now(),
		},
	}

	// Update the listing
	_, err = listingCollection.UpdateOne(ctx, bson.M{"_id": listingID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour du listing"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Listing mis à jour avec succès"})
}
