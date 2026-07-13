package main

import (
	"context"
	"crypto/tls"
	"goapp/authorize"
	"goapp/db"
	"goapp/facility"
	"goapp/myRedis"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
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
	redisHost := os.Getenv("REDIS_ADDR")
	if redisHost == "" {
		redisHost = "localhost:3003"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		redisPassword = "password0"
	}

	secretKey := os.Getenv("SESSION_SECRET")
	if secretKey == "" {
		slog.Error("can not found secret key app is down")
		os.Exit(1)
	}
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	mysql.RegisterTLSConfig("custom-tls", tlsConfig)
	mysqlConnection, err := sqlx.Open("mysql", connectionString)
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		os.Exit(1)
	}
	defer mysqlConnection.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := mysqlConnection.PingContext(ctx); err != nil {
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

	myPool := &redigo.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redigo.Conn, error) {
			return redigo.Dial("tcp", redisHost, redigo.DialPassword(redisPassword), redigo.DialUseTLS(true), redigo.DialConnectTimeout(3*time.Second))
		},
	}
	mySes := myRedis.NewPrefixRedisStoreWithPool(myPool, "hotelService", 3600, []byte(secretKey))

	userDB := db.NewUserDB(mysqlConnection)
	authHandler := authorize.NewAuthHandlerInterface(publicKey, userDB, mySes)
	newApp.Use(authorize.CheckAccessToken(authHandler))
	facilityRepository := facility.NewFacilityRepository(userDB)
	facilityHandler := facility.NewFacilityHandler(facilityRepository)
	newApp.GET("/facility", facilityHandler.Get)
	newApp.Run(":" + port) // listen and serve on
}
