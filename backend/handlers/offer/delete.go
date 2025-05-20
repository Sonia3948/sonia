
package offer

import (
	"context"
	"net/http"
	"time"

	"elsaidaliya/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Delete deletes an offer
func Delete(c *gin.Context) {
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

	// Only allow deletion of user's own offers (or admin can delete any)
	userRole, _ := c.Get("userRole")
	filter := bson.M{"_id": offerID}
	if userRole != "admin" {
		filter["supplierID"] = userID
	}

	result, err := offerCollection.DeleteOne(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression de l'offre"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Offre non trouvée ou vous n'avez pas l'autorisation de la supprimer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Offre supprimée avec succès"})
}
