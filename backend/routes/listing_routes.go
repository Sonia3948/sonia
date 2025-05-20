
package routes

import (
	"elsaidaliya/handlers"
	"elsaidaliya/middleware"

	"github.com/gin-gonic/gin"
)

// SetupListingRoutes configures routes for listing management
func SetupListingRoutes(r *gin.Engine) {
	listings := r.Group("/api/listings")
	{
		listings.GET("/", handlers.GetAllListings)
		listings.GET("/:id", handlers.GetListingByID)
		listings.POST("/", middleware.RequireAuth, middleware.RequireRole("fournisseur"), handlers.CreateListing)
		listings.PUT("/:id", middleware.RequireAuth, middleware.RequireRole("fournisseur"), handlers.UpdateListing)
		listings.DELETE("/:id", middleware.RequireAuth, handlers.DeleteListing)
		listings.GET("/search", handlers.SearchListings)
	}
}
