package handlers

import (
	"elsaidaliya/handlers/auth"

	"go.mongodb.org/mongo-driver/mongo"
)

// InitAllHandlers centralise l'initialisation de tous les sous‑packages handlers.
func InitAllHandlers(db *mongo.Database) {
	auth.Init(db)
	// notification.Init(db) // Commented out since notification.Init is undefined
	// Ajoute ici d'autres modules (ex: user.Init(db), offer.Init(db), etc.)
}
