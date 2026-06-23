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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port if not specified
	}
	newApp.Run(":" + port) // listen and serve on
}
