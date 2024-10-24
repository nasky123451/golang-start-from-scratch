package chat

import (
	"example.com/m/chat/config"
	"example.com/m/chat/handlers"
	"github.com/gin-gonic/gin"
)

func ChatServer() {
	// Initialize configurations, databases, and other services
	config.Init()

	r := gin.Default()

	// Setup routes
	handlers.SetupRoutes(r)

	// Start the server
	r.Run(":8080")
}
