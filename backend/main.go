package main

import (
	"log"

	"elsaidaliya/config"
	"elsaidaliya/database"
	"elsaidaliya/handlers"
	"elsaidaliya/middleware"
	"elsaidaliya/routes"
	"elsaidaliya/server"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	_, db, err := database.Connect(cfg.MongoURI)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize handlers with database connection
	handlers.InitAllHandlers(db)

	// Initialize middleware
	middleware.InitAuthMiddleware(db)

	// Create and configure server
	srv := server.New(db)

	// Set up routes
	srv.SetupBasicRoutes()
	routes.Setup(srv.Router())

	// Start the server
	log.Fatal(srv.Start(cfg.Port))
}
