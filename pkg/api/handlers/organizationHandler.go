package handlers

import (
	"assessment/pkg/database/mongodb/models"
	"assessment/pkg/database/mongodb/repository"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateOrganization creates a new organization record.
func CreateOrganization(c *gin.Context) {
	var org models.Organization
	repo := repository.NewOrganizationRepo()

	err := c.ShouldBindJSON(&org)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}
	orgID, err := repo.CreateOrganization(&org)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create organization"})
		return
	}
	// Return a success message
	c.JSON(http.StatusCreated, gin.H{"organization_id": orgID})
}

// GetOrganizationById retrieves an organization by its ID.
func GetOrganizationById(c *gin.Context) {
	organizationID := c.Param("organization_id")

	repo := repository.NewOrganizationRepo()
	organization, err := repo.GetOrganizationById(organizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch organization details"})
		return
	}

	// Return a success message
	c.JSON(http.StatusOK, models.Organization{
		Id:          organization.Id,
		Name:        organization.Name,
		Description: organization.Description,
	})
}

// GetAllOrganizations lists all organizations.
func GetAllOrganizations(c *gin.Context) {
	repo := repository.NewOrganizationRepo()
	organizations, err := repo.GetAllOrganizations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch organizations"})
		return
	}

	// Return a success message
	c.JSON(http.StatusOK, organizations)
}

// UpdateOrganization updates an existing organization's details.
func UpdateOrganization(c *gin.Context) {
	organizationID := c.Param("organization_id")

	var updateData models.OrganizationUpdate
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("req body", updateData)
	repo := repository.NewOrganizationRepo()

	organization, err := repo.UpdateOrganization(organizationID, &updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update organization"})
		return
	}

	// Return a success message
	c.JSON(http.StatusOK, gin.H{
		"organization_id": organization.Id,
		"name":            organization.Name,
		"description":     organization.Description,
	})
}

// DeleteOrganization removes an organization from the database.
func DeleteOrganization(c *gin.Context) {
	organizationID := c.Param("organization_id")

	repo := repository.NewOrganizationRepo()

	err := repo.DeleteOrganization(organizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete organization"})
		return
	}
	// Return a success message
	c.JSON(http.StatusOK, gin.H{"message": "Organization deleted successfully"})
}

// InviteUserToOrganization sends an invitation to join an organization.
func InviteUserToOrganization(c *gin.Context) {
	organizationID := c.Param("organization_id")
	var requestBody models.InviterequestBody

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	repo := repository.NewOrganizationRepo()
	// Call the InviteUserToOrganization method in the repository
	err := repo.InviteUserToOrganization(organizationID, requestBody.UserEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to invite user to organization"})
		return
	}

	// Return a success message
	c.JSON(http.StatusOK, gin.H{"message": "User invited to organization"})
}
