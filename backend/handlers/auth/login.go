
package auth

import (
	"context"
	"net/http"
	"time"

	"elsaidaliya/models"
	
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// LoginUser function handles both regular user and admin login
func LoginUser(c *gin.Context) {
	var credentials struct {
		Identifier string `json:"identifier"` // Email, phone, or username
		Password   string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// First check if it's the admin trying to login
	var admin models.AdminUser
	err := adminCollection.FindOne(ctx, bson.M{"phone": credentials.Identifier}).Decode(&admin)
	if err == nil {
		// Found admin user, verify password
		err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(credentials.Password))
		if err == nil {
			// Password is correct, login successful
			c.JSON(http.StatusOK, gin.H{
				"message": "Login successful",
				"user": gin.H{
					"id":   admin.ID.Hex(),
					"role": "admin",
				},
			})
			return
		}
	}

	// Not an admin or wrong admin password, try regular user login
	var user models.User
	filter := bson.M{
		"$or": []bson.M{
			{"email": credentials.Identifier},
			{"phone": credentials.Identifier},
		},
	}

	err = userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Identifiant ou mot de passe incorrect"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la connexion"})
		}
		return
	}

	// Check if account is active
	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "Votre compte est en attente d'activation"})
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Identifiant ou mot de passe incorrect"})
		return
	}

	// All good, login successful
	c.JSON(http.StatusOK, gin.H{
		"message": "Connexion réussie",
		"user": gin.H{
			"id":      user.ID.Hex(),
			"name":    user.BusinessName,
			"role":    user.Role,
			"email":   user.Email,
			"phone":   user.Phone,
		},
	})
}
