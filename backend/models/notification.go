
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NotificationType defines the type of notification
type NotificationType string

const (
	NotificationTypeListing NotificationType = "listing"
	NotificationTypeOffer   NotificationType = "offer"
	NotificationTypePayment NotificationType = "payment"
)

// Notification represents a notification in the system
type Notification struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID      primitive.ObjectID `bson:"userID" json:"userID"`
	FromID      primitive.ObjectID `bson:"fromID" json:"fromID"`
	FromName    string             `bson:"fromName" json:"fromName"`
	Type        NotificationType   `bson:"type" json:"type"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Read        bool               `bson:"read" json:"read"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}
