
package offer

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"elsaidaliya/models"
)

// GetAll fetches all offers
func GetAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Define options to sort by createdAt DESC
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"createdAt": -1})

	// Get filter from query params
	filter := bson.M{}
	if supplierID := c.Query("supplier"); supplierID != "" {
		objectID, err := primitive.ObjectIDFromHex(supplierID)
		if err == nil {
			filter["supplierID"] = objectID
		}
	}

	cursor, err := offerCollection.Find(ctx, filter, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des offres"})
		return
	}
	defer cursor.Close(ctx)

	var offers []models.Offer
	if err = cursor.All(ctx, &offers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du traitement des offres"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"offers": offers})
}

// GetByID fetches a single offer by ID
func GetByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	offerID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID d'offre invalide"})
		return
	}

	var offer models.Offer
	err = offerCollection.FindOne(ctx, bson.M{"_id": offerID}).Decode(&offer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Offre non trouvée"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération de l'offre"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"offer": offer})
}
