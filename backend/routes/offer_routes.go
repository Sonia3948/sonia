
package routes

import (
	"elsaidaliya/handlers"
	"elsaidaliya/middleware"

	"github.com/gin-gonic/gin"
)

// SetupOfferRoutes configures routes for offer management
func SetupOfferRoutes(r *gin.Engine) {
	offers := r.Group("/api/offers")
	{
		offers.GET("/", handlers.GetAllOffers)
		offers.GET("/:id", handlers.GetOfferByID)
		offers.POST("/", middleware.RequireAuth, middleware.RequireRole("fournisseur"), handlers.CreateOffer)
		offers.PUT("/:id", middleware.RequireAuth, middleware.RequireRole("fournisseur"), handlers.UpdateOffer)
		offers.DELETE("/:id", middleware.RequireAuth, handlers.DeleteOffer)
	}
}
