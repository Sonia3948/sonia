
package auth

import (
	"context"
	"net/http"
	"time"

	"elsaidaliya/models"
	"elsaidaliya/utils"
	
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ForgotPassword handles password reset requests
func ForgotPassword(c *gin.Context) {
	var request struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Adresse e-mail invalide"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"email": request.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Don't reveal if email exists or not for security
			c.JSON(http.StatusOK, gin.H{"message": "Si votre email est enregistré, vous recevrez un lien de réinitialisation"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la vérification de l'email"})
		return
	}

	// Generate reset token
	resetToken := utils.GenerateRandomToken()
	resetExpiry := time.Now().Add(24 * time.Hour)

	// Update user with reset token
	_, err = userCollection.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		bson.M{
			"$set": bson.M{
				"resetPasswordToken": resetToken,
				"resetPasswordExpires": resetExpiry,
				"updatedAt": time.Now(),
			},
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la génération du token de réinitialisation"})
		return
	}

	// TODO: Send email with reset token
	// For now, just return success
	c.JSON(http.StatusOK, gin.H{"message": "Si votre email est enregistré, vous recevrez un lien de réinitialisation"})
}
