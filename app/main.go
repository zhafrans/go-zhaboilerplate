package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"

	"go-zhaboilerplate/auth"
	mysqlRepo "go-zhaboilerplate/internal/repository/mysql"
	"go-zhaboilerplate/internal/rest"
	"go-zhaboilerplate/internal/rest/middleware"
	"go-zhaboilerplate/user"
)

const (
	defaultTimeout   = 30
	defaultAddress   = ":9090"
	defaultSecret    = "boilerplate-secret-key-change-me"
	defaultJWTExpiry = 24
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Notice: .env file not found or failed to load, falling back to system environment variables")
	}
}

func main() {
	// Prepare database
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USER")
	dbPass := os.Getenv("DATABASE_PASS")
	dbName := os.Getenv("DATABASE_NAME")
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(`mysql`, dsn)
	if err != nil {
		log.Fatal("failed to open connection to database: ", err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal("failed to ping database: ", err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal("got error when closing the DB connection: ", err)
		}
	}()

	// Prepare JWT configuration
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = defaultSecret
	}

	jwtExpiryHours, err := strconv.Atoi(os.Getenv("JWT_EXPIRY_HOURS"))
	if err != nil || jwtExpiryHours <= 0 {
		jwtExpiryHours = defaultJWTExpiry
	}
	jwtExpiry := time.Duration(jwtExpiryHours) * time.Hour

	// Prepare Fiber application
	app := fiber.New(fiber.Config{
		AppName: "go-zhaboilerplate",
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(middleware.CORS)

	timeoutStr := os.Getenv("CONTEXT_TIMEOUT")
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		log.Println("failed to parse timeout, using default timeout")
		timeout = defaultTimeout
	}
	timeoutContext := time.Duration(timeout) * time.Second
	app.Use(middleware.SetRequestContextWithTimeout(timeoutContext))

	// Prepare Repository
	userRepo := mysqlRepo.NewUserRepository(dbConn)

	// Prepare Service layer
	userSvc := user.NewService(userRepo)
	authSvc := auth.NewService(userRepo, jwtSecret, jwtExpiry)

	// Middleware
	authMiddleware := middleware.JWTMiddleware(jwtSecret)

	// Register REST handlers
	rest.NewAuthHandler(app, authSvc, authMiddleware)
	rest.NewUserHandler(app, userSvc, authMiddleware)

	// Start Server
	address := os.Getenv("SERVER_ADDRESS")
	if address == "" {
		address = defaultAddress
	}
	log.Fatal(app.Listen(address))
}
