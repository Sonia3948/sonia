
package offer

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var offerCollection *mongo.Collection

// InitHandlers initializes handlers for offer management
func InitHandlers(db *mongo.Database) {
	offerCollection = db.Collection("offers")
}
