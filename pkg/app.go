package pkg

import (
	"assessment/pkg/api/routes"
	db "assessment/pkg/database"

	"github.com/gin-gonic/gin"
)

// Run starts the web server and registers the API routes.
func Run() {
	// Initialize the Gin router with default middleware.
	router := gin.Default()

	// Connect to the database.
	err := db.Connect()
	if err != nil {
		panic(err)
	}

	// Register the API routes with the router.
	routes.RegisterRoutes(router)

	// Start the web server on port  8080.
	router.Run(":8080")
}
