
package routes

import (
	"github.com/gin-gonic/gin"
)

// Setup initializes all routes
func Setup(r *gin.Engine) {
	SetupAuthRoutes(r)
	SetupUserRoutes(r)
	SetupListingRoutes(r)
	SetupOfferRoutes(r)
	SetupNotificationRoutes(r)
}
