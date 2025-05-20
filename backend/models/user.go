
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system
type User struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	BusinessName     string             `bson:"businessName" json:"businessName"`
	Role             string             `bson:"role" json:"role"` // admin, pharmacien, fournisseur
	Phone            string             `bson:"phone" json:"phone"`
	Email            string             `bson:"email" json:"email"`
	Password         string             `bson:"password" json:"-"` // Password hash, not returned in JSON
	Wilaya           string             `bson:"wilaya" json:"wilaya"`
	RegisterImageURL string             `bson:"registerImageUrl" json:"registerImageUrl"`
	RegisterNumber   string             `bson:"registerNumber" json:"registerNumber"` // Registre de commerce
	IsActive         bool               `bson:"isActive" json:"isActive"`
	Subscription     string             `bson:"subscription" json:"subscription"` // bronze, argent, or
	SubExpiry        time.Time          `bson:"subExpiry" json:"subExpiry"`
	CreatedAt        time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt        time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// AdminUser represents an admin user
type AdminUser struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Phone     string             `bson:"phone" json:"phone"`
	Password  string             `bson:"password" json:"-"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}
