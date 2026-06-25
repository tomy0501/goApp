package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
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
	sessionSecret := os.Getenv("SESSION_SECRET")
	if redisAddr == "" {
		redisAddr = "localhost:6379" // default Redis address if not specified
	}
	if redisPassword == "" {
		redisPassword = "" // default Redis password if not specified
	}
	if sessionSecret == "" {
		sessionSecret = "secret" // default session secret if not specified
	}

	store, err := redis.NewStore(10, "tcp", redisAddr, redisPassword, sessionSecret)
	if err != nil {
		log.Printf("Error creating Redis store: %v", err)
		os.Exit(1)
	}

	log.Printf("Connected to Redis at %s", redisAddr)
	newApp.Use(sessions.Sessions("mysession", store))
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
