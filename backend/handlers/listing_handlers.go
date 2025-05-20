
package handlers

import (
	"elsaidaliya/handlers/listing"
	"go.mongodb.org/mongo-driver/mongo"
)

// InitListingHandlers initializes handlers for listing management
func InitListingHandlers(db *mongo.Database) {
	listing.InitListingHandlers(db)
}

// GetAllListings fetches all medicine listings
var GetAllListings = listing.GetAllListings

// GetListingByID fetches a single listing by ID
var GetListingByID = listing.GetListingByID

// CreateListing creates a new medicine listing
var CreateListing = listing.CreateListing

// UpdateListing updates an existing medicine listing
var UpdateListing = listing.UpdateListing

// DeleteListing deletes a medicine listing
var DeleteListing = listing.DeleteListing

// SearchListings searches for medicines in listings
var SearchListings = listing.SearchListings
