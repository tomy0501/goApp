package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"goapp/authorize"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
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

	pemString, err := os.ReadFile("public.pem")
	if err != nil {
		log.Fatalf("Error reading public key: %v", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pemString)
	if err != nil {
		log.Fatalf("Error parsing public key: %v", err)
	}

	authhandler := authorize.AuthHandler{
		PublicKey: publicKey,
	}

	newApp.Use(authorize.CheckAccessToken(&authhandler))

	newApp.GET("/facility", func(c *gin.Context) {

	})
	newApp.Run(":" + port) // listen and serve on
}
