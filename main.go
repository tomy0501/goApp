package main

import (
	"crypto/tls"
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
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
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisAddr == "" {
		redisAddr = "localhost:6379" // default Redis address if not specified
	}
	if redisPassword == "" {
		redisPassword = "" // default Redis password if not specified
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	})
	_, err = redisClient.Ping().Result()
	if err != nil {
		log.Printf("Error connecting to Redis: %v", err)
		os.Exit(1)
	}
	defer redisClient.Close()
	redisClient.Set("test_key", "test_value", 0)
	log.Printf("Connected to Redis at %s", redisAddr)
	connectionString := os.Getenv("DB_CONNECTION_STRING")
	if connectionString == "" {
		connectionString = "user:password@tcp(localhost:3306)/dbname" // default connection string if not specified
	}
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		os.Exit(1)
	}
	defer db.Close()
	newApp.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	newApp.Run(":" + port) // listen and serve on
}
