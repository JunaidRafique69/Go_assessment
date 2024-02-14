package routes

import (
	"assessment/pkg/api/handlers"
	"assessment/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up the application's HTTP routes.
func RegisterRoutes(router *gin.Engine) {

	// Define authentication routes.
	auth := router.Group("/auth")
	{
		auth.POST("/signup", handlers.Signup)                           // User registration
		auth.POST("/signin", handlers.SignIn)                           // User login
		auth.POST("/refresh-token", handlers.RefreshToken)              // Token refresh
		auth.POST("/revoke-refresh-token", handlers.RevokeRefreshToken) // Token revocation
	}

	// Define organization routes, secured with authentication.
	organization := router.Group("/api")
	organization.Use(middleware.AuthMiddleware())
	{
		organization.POST("organization", handlers.CreateOrganization)                                                  // Organization creation
		organization.GET("/organization/:organization_id", middleware.InviteMiddleware(), handlers.GetOrganizationById) // Organization retrieval with invitation check
		organization.GET("/organization", handlers.GetAllOrganizations)                                                 // All organizations retrieval
		organization.PUT("/organization/:organization_id", handlers.UpdateOrganization)                                 // Organization update
		organization.DELETE("/organization/:organization_id", handlers.DeleteOrganization)                              // Organization deletion
		organization.POST("/organization/:organization_id/invite", handlers.InviteUserToOrganization)                   // Organization invitation
	}
}
