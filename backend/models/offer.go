
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Offer represents a promotional offer from a supplier
type Offer struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	SupplierID  primitive.ObjectID `bson:"supplierID" json:"supplierID"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Price       string             `bson:"price" json:"price"`
	ImageURL    string             `bson:"imageUrl" json:"imageUrl"`
	ExpiresAt   time.Time          `bson:"expiresAt" json:"expiresAt"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}
