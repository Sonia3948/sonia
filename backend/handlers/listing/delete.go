
package listing

import (
	"context"
	"net/http"
	"time"

	"elsaidaliya/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Delete deletes a medicine listing
func Delete(c *gin.Context) {
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

	// Only allow deletion of user's own listings (or admin can delete any)
	userRole, _ := c.Get("userRole")
	filter := bson.M{"_id": listingID}
	if userRole != "admin" {
		filter["supplierID"] = userID
	}

	result, err := listingCollection.DeleteOne(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression du listing"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Listing non trouvé ou vous n'avez pas l'autorisation de le supprimer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Listing supprimé avec succès"})
}
