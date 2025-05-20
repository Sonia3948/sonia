
package listing

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var listingCollection *mongo.Collection

// InitHandlers initializes handlers for listing management
func InitHandlers(db *mongo.Database) {
	listingCollection = db.Collection("listings")
}
