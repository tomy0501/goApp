package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	newApp := gin.Default()
	newApp.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins != "" {
		newApp.Use(func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", allowedOrigins)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Next()
		})
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port if not specified
	}
	newApp.Run(":" + port) // listen and serve on
}
