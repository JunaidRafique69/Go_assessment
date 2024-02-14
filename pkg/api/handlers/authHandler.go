package handlers

import (
	"assessment/config"
	"assessment/pkg/database/mongodb/models"
	"assessment/pkg/database/mongodb/repository"
	"assessment/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Signup handles the creation of a new user account.
func Signup(c *gin.Context) {
	// Parse and validate the incoming JSON payload.
	var user models.User
	repo := repository.NewUserRepo()

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	// Validate the user data before proceeding.
	if err := utils.ValidateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the user's password for secure storage.
	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = hash

	// Attempt to create the user in the database.
	createdUser, err := repo.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not created"})
		return
	}

	// Generate authentication tokens for the newly created user.
	access_token, refresh_token, err := utils.GenerateTokens(createdUser.Name, createdUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Respond with success message and tokens.
	c.JSON(http.StatusCreated, models.AuthResponse{
		Message:      "User created successfully",
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	})
}

// SignIn authenticates a user based on their credentials.
func SignIn(c *gin.Context) {
	// Parse the incoming JSON payload containing user credentials.
	var credentials *models.Credentials
	repo := repository.NewUserRepo()

	err := c.ShouldBindJSON(&credentials)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	// Check that both username and password are provided.
	if credentials.Email == "" || credentials.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password cannot be empty"})
		return
	}

	// Find the user by email in the database.
	userFound, err := repo.FindUserByEmail(credentials.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Verify the provided password against the stored hash.
	isMatch, err := utils.CheckPasswordHash(credentials.Password, userFound.Password)
	if err != nil || !isMatch {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate authentication tokens for the authenticated user.
	access_token, refresh_token, err := utils.GenerateTokens(userFound.Name, userFound.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Respond with success message and tokens.
	c.JSON(http.StatusOK, models.AuthResponse{
		Message:      "SignIn successful",
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	})
}

// RefreshToken issues new access and refresh tokens using a valid refresh token.
func RefreshToken(c *gin.Context) {
	// Parse the incoming JSON payload containing the refresh token.
	var request models.RefreshToken

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	// Verify the refresh token and extract the associated username and email.
	username, email, err := utils.VerifyRefreshToken(request.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Generate new access and refresh tokens for the user.
	accessToken, refreshToken, err := utils.GenerateTokens(username, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	// Prepare the response with the new tokens.
	response := models.AuthResponse{
		Message:      "Tokens refreshed successfully",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, response)
}

// RevokeRefreshToken removes a refresh token from the system, effectively logging the user out.
func RevokeRefreshToken(c *gin.Context) {
	// Parse the incoming JSON payload containing the refresh token to be revoked.
	var requestBody *models.RefreshToken

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	// Delete the refresh token from Redis
	err := config.Init_redis().Del(requestBody.Token).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Refresh token revoked"})
}
