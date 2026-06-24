package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("No .env file found, proceeding with environment variables")
		} else {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}
	newApp := gin.Default()

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
	connectionString := os.Getenv("DB_CONNECTION_STRING")
	if connectionString == "" {
		connectionString = "user:password@tcp(localhost:3306)/dbname" // default connection string if not specified
	}
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal("Error connecting to database: %v", err)
	}
	defer db.Close()

	newApp.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	newApp.Run(":" + port) // listen and serve on
}
