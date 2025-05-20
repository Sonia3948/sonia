
package handlers

import (
	"elsaidaliya/handlers/offer"
	"go.mongodb.org/mongo-driver/mongo"
)

// InitOfferHandlers initializes handlers for offer management
func InitOfferHandlers(db *mongo.Database) {
	offer.InitOfferHandlers(db)
}

// GetAllOffers fetches all offers
var GetAllOffers = offer.GetAllOffers

// GetOfferByID fetches a single offer by ID
var GetOfferByID = offer.GetOfferByID

// CreateOffer creates a new offer
var CreateOffer = offer.CreateOffer

// UpdateOffer updates an existing offer
var UpdateOffer = offer.UpdateOffer

// DeleteOffer deletes an offer
var DeleteOffer = offer.DeleteOffer
