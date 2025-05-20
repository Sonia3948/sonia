
package listing

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

// GetAll fetches all medicine listings
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

	cursor, err := listingCollection.Find(ctx, filter, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des listings"})
		return
	}
	defer cursor.Close(ctx)

	var listings []models.Listing
	if err = cursor.All(ctx, &listings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du traitement des listings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"listings": listings})
}

// GetByID fetches a single listing by ID
func GetByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	listingID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de listing invalide"})
		return
	}

	var listing models.Listing
	err = listingCollection.FindOne(ctx, bson.M{"_id": listingID}).Decode(&listing)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Listing non trouvé"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération du listing"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"listing": listing})
}

// Search searches for medicines in listings
func Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Paramètre de recherche manquant"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build the search query - search in medication names and in listing title/description
	filter := bson.M{
		"$or": []bson.M{
			{"title": bson.M{"$regex": query, "$options": "i"}},
			{"description": bson.M{"$regex": query, "$options": "i"}},
			{"medications.name": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	cursor, err := listingCollection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la recherche"})
		return
	}
	defer cursor.Close(ctx)

	var results []models.Listing
	if err = cursor.All(ctx, &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du traitement des résultats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"listings": results})
}
