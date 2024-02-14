package middleware

import (
	"assessment/pkg/database/mongodb/repository"
	"assessment/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware checks for a valid authorization token in the request headers.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the Authorization header from the request.
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		// Extract the token string after removing the "Bearer" prefix.
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the extracted token.
		_, err := utils.ValidateToken(tokenString)
		if err != nil {
			// If the token is invalid, respond with an Unauthorized status.
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Proceed to the next handler if the token is valid.
		c.Next()
	}
}

// InviteMiddleware verifies if the user is authorized to perform actions related to invitations.
func InviteMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the organization ID from the URL parameter.
		organizationID := c.Param("organization_id")
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		// Extract the token string after removing the "Bearer" prefix.
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		userEmail, err := utils.GetEmailFromToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Create a new instance of the Organization repository.
		repo := repository.NewOrganizationRepo()
		organization, err := repo.GetOrganizationById(organizationID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get organization"})
			c.Abort()
			return
		}

		// Check if the user is in the list of invited users for the organization.
		isInvited := false
		for _, invitedUser := range organization.InvitedUsers {
			if invitedUser == userEmail {
				isInvited = true
				break
			}
		}

		// If the user is not invited, respond with an Unauthorized status.
		if !isInvited {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User is not invited to the organization"})
			c.Abort()
			return
		}

		c.Next()
	}
}
