package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

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
	pool := &redigo.Pool{
		Dial: func() (redigo.Conn, error) {
			return redigo.Dial("tcp", redisAddr, redigo.DialUseTLS(true), redigo.DialPassword(redisPassword))
		},
	}
	store, err := redis.NewStoreWithPool(pool, []byte(sessionSecret))
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
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	mysql.RegisterTLSConfig("custom-tls", tlsConfig)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Printf("Error pinging database: %v", err)
		os.Exit(1)
	}
	newApp.GET("/ping", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("ping", "pong")
		session.Save()
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	newApp.GET("/user", func(c *gin.Context) {
		getUserSQL := "SELECT user_id, username, email FROM user WHERE user_id = ? LIMIT 1"
		result, err := db.Query(getUserSQL, 1) // Example: Fetch user with ID 1
		if err != nil {
			log.Printf("Error querying database: %v", err)
			c.JSON(500, gin.H{"error": "Internal Server Error"})
			return
		}
		defer result.Close()
		var user User
		if result.Next() {
			err := result.Scan(&user.ID, &user.Username, &user.Email)
			if err != nil {
				log.Printf("Error scanning user data: %v", err)
				c.JSON(500, gin.H{"error": "Internal Server Error"})
				return
			}
		}
		c.JSON(200, user)
	})
	newApp.Run(":" + port) // listen and serve on
}
