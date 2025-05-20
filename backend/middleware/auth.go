package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"elsaidaliya/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection
var adminCollection *mongo.Collection

// InitAuthMiddleware initializes the collections for auth middleware
func InitAuthMiddleware(db *mongo.Database) {
	userCollection = db.Collection("users")
	adminCollection = db.Collection("admins")
}

// RequireAuth middleware to check if user is authenticated
func RequireAuth(c *gin.Context) {
	// Get token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentification requise"})
		c.Abort()
		return
	}

	// Extract token (remove "Bearer " prefix)
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Format d'authentification invalide"})
		c.Abort()
		return
	}

	// Verify token and get user ID
	// TODO: Implement proper JWT verification
	// For now, we'll just check if the token exists in the user collection
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"sessionToken": tokenString}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Check if it's an admin token
			var admin models.AdminUser
			err = adminCollection.FindOne(ctx, bson.M{"sessionToken": tokenString}).Decode(&admin)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalide ou expiré"})
				c.Abort()
				return
			}
			
			// It's an admin
			c.Set("userID", admin.ID.Hex())
			c.Set("userRole", "admin")
			c.Next()
			return
		}
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur de vérification d'authentification"})
		c.Abort()
		return
	}
	
	// User is authenticated
	c.Set("userID", user.ID.Hex())
	c.Set("userRole", user.Role)
	c.Next()
}

// RequireRole middleware to check if user has required role
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentification requise"})
			c.Abort()
			return
		}
		
		role := userRole.(string)
		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}
		
		c.JSON(http.StatusForbidden, gin.H{"error": "Accès non autorisé"})
		c.Abort()
	}
}

// GetUserID extracts the userID from the context
func GetUserID(c *gin.Context) (primitive.ObjectID, error) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		return primitive.ObjectID{}, fmt.Errorf("user ID not found in context")
	}
	
	return primitive.ObjectIDFromHex(userIDStr.(string))
}
